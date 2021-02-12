package apc

import (
	"html/template"

	"github.com/golangcollege/sessions"
)

// AppContext holds the local app context for the application.
type AppContext struct {
	Session       *sessions.Session
	TemplateCache map[string]*template.Template
}

// New creates a new Apc struct.
func New(session *sessions.Session) *AppContext {
	return &AppContext{
		Session:       session,
		TemplateCache: nil,
	}
}
