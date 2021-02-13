package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	td "github.com/sophiabrandt/go-party-finder/business/data"
	"github.com/sophiabrandt/go-party-finder/business/data/party"
	"github.com/sophiabrandt/go-party-finder/business/app/forms"
	"github.com/sophiabrandt/go-party-finder/foundation/web"
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

	parties, err := pg.party.Query(ctx, v.TraceID, pageNumber, rowsPerPage)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, "home.page.tmpl", &td.TemplateData{Parties: parties}, http.StatusOK)
}

// querybyID shows the details page for a given party.
func (pg partyGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	id := web.Param(r, "id")
	prty, err := pg.party.QueryByID(ctx, v.TraceID, id)
	if err != nil {
		switch errors.Cause(err) {
		case party.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case party.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", id)
		}
	}

	flash := v.AppContext.Session.PopString(r, "flash")
	return web.Respond(ctx, w, "party_detail.page.tmpl", &td.TemplateData{Party: prty, Flash: flash}, http.StatusOK)
}

// createForm shows the web form for creating a new party.
func (pg partyGroup) createForm(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	form := forms.New(nil)
	return web.Respond(ctx, w, "create.page.tmpl", &td.TemplateData{Form: form}, http.StatusOK)
}

// create parses form data and creates a new party.
func (pg partyGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var np party.NewParty
	if form, err := web.DecodeForm(r, &np); err != nil {
		return web.Respond(ctx, w, "create.page.tmpl", &td.TemplateData{Form: form}, http.StatusUnprocessableEntity)
	}

	prty, err := pg.party.Create(ctx, v.TraceID, np, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new party: %+v", np)
	}

	v.AppContext.Session.Put(r, "flash", "Party successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/parties/%v", prty.ID), http.StatusSeeOther)
	return nil
}
