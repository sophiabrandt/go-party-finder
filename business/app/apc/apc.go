package apc

import (
	"github.com/golangcollege/sessions"
	"html/template"
)

// AppContext holds the local app context for the application.
type AppContext struct {
	Session       *sessions.Session
	TemplateCache map[string]*template.Template
}

// New creates a new Apc struct.
func New(session *sessions.Session, tc map[string]*template.Template) *AppContext {
	return &AppContext{
		Session:       session,
		TemplateCache: tc,
	}
}
