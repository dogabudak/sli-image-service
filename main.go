package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/images", func(r chi.Router) {
		r.Route("/{imageID}", func(r chi.Router) {
			r.Use(ImageCtx)
			r.Get("/", getImage)
		})
	})
	http.ListenAndServe(":6057", r)
}

func ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageID := chi.URLParam(r, "imageID")
		image, err := dbGetImage(imageID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "image", image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	image, ok := ctx.Value("image").(*Image)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	w.Write([]byte(fmt.Sprintf("title:%s", image.Name)))
}

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var images = []*Image{
	{ID: "1", Name: "Example.jpeg"},
}

func dbGetImage(id string) (*Image, error) {
	for _, a := range images {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("article not found.")
}
