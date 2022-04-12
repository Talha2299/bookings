package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Talha2299/bookings/internal/config"
	"github.com/Talha2299/bookings/internal/handlers"
	"github.com/Talha2299/bookings/internal/helpers"
	"github.com/Talha2299/bookings/internal/models"
	"github.com/Talha2299/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNum = ":8080"

var App config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is main application function
func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Application is starting on port %s", portNum))

	srv := &http.Server{
		Addr:    portNum,
		Handler: Routes(&App),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	//why am I going to put into the session
	gob.Register(models.Reservation{})

	//change to true when in production
	App.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	App.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	App.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = App.InProduction

	App.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	App.TemplateCache = tc
	App.UseCache = false

	repo := handlers.NewRepo(&App)
	handlers.NewHandler(repo)

	render.NewTemplates(&App)
	helpers.NewHelpers(&App)

	return nil
}
