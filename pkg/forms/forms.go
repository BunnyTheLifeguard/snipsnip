package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form struct embeds a url.Values object to hold form data & validation errors
type Form struct {
	url.Values
	Errors errors
}

// New initializes a custom Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if fields are blank
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// MaxLength checks the number of characters in a field
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is limited to %d characters.", d))
	}
}

// PermittedValues checks if a specific field matches a set of permitted values
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid.")
}

// Valid returns true if no errors present
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
