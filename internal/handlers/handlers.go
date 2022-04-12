package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Talha2299/bookings/internal/config"
	"github.com/Talha2299/bookings/internal/forms"
	"github.com/Talha2299/bookings/internal/helpers"
	"github.com/Talha2299/bookings/internal/models"
	"github.com/Talha2299/bookings/internal/render"
)

// Repo  the repository use by handlers
var Repo *Repository

// Repository is a repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandler sets the repository for handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTeamplates(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.RenderTeamplates(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservations renders the make-reservation page and display form
func (m *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTeamplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservations handles the posting of the reservatoin form
func (m *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTeamplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Genrals render the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTeamplates(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors render the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTeamplates(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability render the search Availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTeamplates(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability render the search Availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start date is %s and End date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handle request for Availability and respond back in JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact render the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTeamplates(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary render the contact page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTeamplates(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
