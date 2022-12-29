package dbrepo

import (
	"errors"
	"time"

	"github.com/eldicela/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservtion Inserts a reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

// func (m *testDBRepo) GetLastInsertedID() (int, error) {

// 	var newId int

// 	return newId, nil
// }

// InsertRoomRestriction inserts a room restriction into database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID Returns true if availability exist for roomId and false if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// GetRoomById Gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {

	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

// AllReservations Returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations Returns a slice of all reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationById returns one reservation by ID
func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {

	var res models.Reservation

	return res, nil
}

// UpdateReservation updates a user in the database
func (m *testDBRepo) UpdateReservation(u models.Reservation) error {

	return nil
}

// DeleteReservation one reservation by Id
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

// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

// InsertBlockForRoom deletes a room restriction
func (m *testDBRepo) DeleteBlockById(id int) error {

	return nil
}
