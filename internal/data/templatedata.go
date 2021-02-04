package data

import "github.com/sophiabrandt/go-party-finder/internal/data/party"

type TemplateData struct {
	Party   *party.Party
	Parties []*party.Info
}
