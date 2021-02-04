package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	td "github.com/sophiabrandt/go-party-finder/internal/data"
	"github.com/sophiabrandt/go-party-finder/internal/data/party"
	"github.com/sophiabrandt/go-party-finder/internal/web"
)

type partyGroup struct {
	party party.Party
}

// query shows the homepage with a list of available parties.
func (pg partyGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	page := web.Param(r, "page")
	// if we are on base route "/" set default page
	if page == "" {
		page = "1"
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid page format: %s", page), http.StatusBadRequest)
	}

	rows := web.Param(r, "rows")
	// if we are on base route "/" set default rows
	if rows == "" {
		rows = "4"
	}
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return web.NewRequestError(fmt.Errorf("invalid rows format: %s", rows), http.StatusBadRequest)
	}

	log.Println(pageNumber, rowsPerPage)

	parties, err := pg.party.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, "home.page.tmpl", &td.TemplateData{Parties: parties}, http.StatusOK)
}
