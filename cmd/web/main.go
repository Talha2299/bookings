package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Talha2299/bookings/internal/config"
	"github.com/Talha2299/bookings/internal/driver"
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

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(App.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Sprintf("Application is starting on port %s", portNum))

	srv := &http.Server{
		Addr:    portNum,
		Handler: Routes(&App),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	//why am I going to put into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	App.MailChan = mailChan

	//change to true when in production
	App.InProduction = *inProduction
	App.UseCache = *useCache

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

	//Connect to database
	log.Println("Connecting to Database.......")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying.........")
	}

	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	App.TemplateCache = tc

	repo := handlers.NewRepo(&App, db)
	handlers.NewHandler(repo)

	render.NewRenderer(&App)
	helpers.NewHelpers(&App)

	return db, nil
}
