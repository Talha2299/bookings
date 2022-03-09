package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Talha2299/bookings/pkj/config"
	"github.com/Talha2299/bookings/pkj/handlers"
	"github.com/Talha2299/bookings/pkj/render"
)

const portNum = ":8080"

// main is main application function
func main() {

	var App config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	App.TemplateCache = tc
	App.UseCache = false

	repo := handlers.NewRepo(&App)
	handlers.NewHandler(repo)

	render.NewTemplates(&App)

	fmt.Println(fmt.Sprintf("Application is starting on port %s", portNum))

	srv := &http.Server{
		Addr:    portNum,
		Handler: Routes(&App),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
