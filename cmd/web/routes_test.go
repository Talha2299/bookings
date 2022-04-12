package main

import (
	"fmt"
	"testing"

	"github.com/Talha2299/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := Routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing: terst passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but %T", v))
	}
}
