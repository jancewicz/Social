package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jancewicz/social/internal/store"
)

type commentKey string

const commentCtx commentKey = "comment"

func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalSeverError(w, r, err)
	}
}

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
	}
	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalSeverError(w, r, err)
		return
	}
}

type UpdateCommenttPayload struct {
	Content *string `json:"content" validate:"omitempty, max=500"`
}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := getCommentFromCtx(r)

	var payload UpdateCommenttPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Content != nil {
		comment.Content = *payload.Content
	}

	if err := app.store.Comments.Update(r.Context(), comment); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalSeverError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalSeverError(w, r, err)
	}
}

func (app *application) commentContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "comId")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalSeverError(w, r, err)
			return
		}

		ctx := r.Context()

		comment, err := app.store.Comments.GetByComID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalSeverError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, comment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCommentFromCtx(r *http.Request) *store.Comment {
	comment, _ := r.Context().Value(commentCtx).(*store.Comment)
	return comment
}
