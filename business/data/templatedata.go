package data

import (
	"github.com/sophiabrandt/go-party-finder/foundation/web/forms"
	"github.com/sophiabrandt/go-party-finder/business/data/party"
)

type TemplateData struct {
	Party       *party.Info
	Parties     []*party.Info
	CurrentYear int
	Form        *forms.Form
	Flash       string
}
