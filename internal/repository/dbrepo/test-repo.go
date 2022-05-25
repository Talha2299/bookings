package dbrepo

import (
	"errors"
	"time"

	"github.com/Talha2299/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	// if the room ID is 2 then fail, otherwise passeed
	if res.RoomID == 2 {
		return 0, errors.New("Soem error")
	}

	return 1, nil
}

// InsertRoomRestriction insert the room restriciton into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	if r.RoomID == 1000 {
		return errors.New("Soem error")
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID an dfalse if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any , for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

// GetUserByID gets a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	if email == "me@here.ca" {
		return 1, "", nil
	}

	return 0, "", errors.New("some error")
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var res models.Reservation
	return res, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil

}

func (m *testDBRepo) DeleteReservation(id int) error {

	return nil

}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil

}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil

}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction
	return restrictions, nil
}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

func (m *testDBRepo) DeleteBlockByID(id int) error {

	return nil
}
