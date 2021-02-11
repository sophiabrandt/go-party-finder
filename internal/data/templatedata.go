package data

import (
	"github.com/sophiabrandt/go-party-finder/internal/data/party"
	"github.com/sophiabrandt/go-party-finder/internal/forms"
)

type TemplateData struct {
	Party       *party.Info
	Parties     []*party.Info
	CurrentYear int
	Form        *forms.Form
	Flash       string
}
