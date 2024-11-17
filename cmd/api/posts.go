package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jancewicz/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO change after auth
		UserID: 1,
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalSeverError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// Same variable as in routing in api.go
	idParam := chi.URLParam(r, "postId")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalSeverError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalSeverError(w, r, err)
	}
}
