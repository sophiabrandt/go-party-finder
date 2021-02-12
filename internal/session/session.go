package session

import (
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/sophiabrandt/go-party-finder/internal/config"
)

var conf *config.Conf

// NewSession sets up the config for the session manager.
func NewSession(c *config.Conf) {
	conf = c
}

// NewStore creates the session manager.
func NewStore() *sessions.Session {
	session := sessions.New([]byte(conf.Web.SessionSecret))
	session.Lifetime = 12 * time.Hour
	session.Persist = true
	session.SameSite = http.SameSiteLaxMode
	session.Secure = conf.Web.InProduction
	return session
}
