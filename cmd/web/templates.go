package main

import "github.com/BunnyTheLifeguard/snipsnip/pkg/models"

type templateData struct {
	Snip  *models.Snip
	Snips []*models.Snip
}
