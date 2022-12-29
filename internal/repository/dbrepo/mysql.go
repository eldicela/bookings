package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/eldicela/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

// InsertReservtion Inserts a reservation into database
func (m *mysqlDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			values (?, ?, ?, ?, ?, ?, ?, ?, ? );  `

	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	)

	// newId, errr := repository.DatabaseRepo.GetLastInsertedID()
	newId, errr := GetLastInsertedID(m.DB)

	if errr != nil {
		return 0, errr
	}
	if err != nil {
		return 0, err
	}

	return newId, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *mysqlDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id)
			values (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		time.Now(),
		time.Now(),
		r.RestrictionId,
	)

	if err != nil {
		return err
	}

	return nil
}

// func (m *mysqlDBRepo) GetLastInsertedID() (int, error) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	var newId int

// 	rows := m.DB.QueryRowContext(ctx, "select last_insert_id()")
// 	err := rows.Scan(&newId)
// 	if err != nil {
// 		log.Println(err)
// 		return 0, err
// 	}

// 	return newId, nil
// }

func GetLastInsertedID(conn *sql.DB) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int

	rows := conn.QueryRowContext(ctx, "select last_insert_id()")
	err := rows.Scan(&newId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return newId, nil
}

// SearchAvailabilityByDatesByRoomID Returns true if availability exist for roomId and false if no availability
func (m *mysqlDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `
	SELECT count(id)
	FROM room_restrictions
	WHERE
		room_id = ? 
		and ? < end_date and ? > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any, for given date range
func (m *mysqlDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
	SELECT r.id, r.room_name
	FROM rooms r
	WHERE r.id not in 
		(	SELECT rr.room_id from room_restrictions rr where ? < rr.end_date and ? > rr.start_date);
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomById Gets a room by id
func (m *mysqlDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		SELECT id, room_name, created_at, updated_at FROM rooms WHERE id = ?; 
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a user by id
func (m *mysqlDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			 FROM users WHERE id =?			
	`
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (m *mysqlDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users set first_name = ?, last_name = ?, email = ?, access_level =?, updated_at = ? `

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates the user
func (m *mysqlDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")

	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations Returns a slice of all reservations
func (m *mysqlDBRepo) AllReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed , rm.id, rm.room_name
	FROM reservations r
	LEFT JOIN rooms rm on (r.room_id = rm.id)
	ORDER BY r.start_date asc;
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(

			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations Returns a slice of all reservations
func (m *mysqlDBRepo) AllNewReservations() ([]models.Reservation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name
	FROM reservations r
	LEFT JOIN rooms rm on (r.room_id = rm.id)
	WHERE processed =0
	ORDER BY r.start_date asc;
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(

			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomId,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationById returns one reservation by ID
func (m *mysqlDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at,
			 r.updated_at, r.processed, rm.id, rm.room_name
			 FROM reservations r
			 LEFT JOIN rooms rm ON (r.room_id = rm.id)
			 WHERE r.id =?
			 `

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomId,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

// UpdateReservation updates a user in the database
func (m *mysqlDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE reservations set first_name = ?, last_name = ?, email = ?, phone = ?, updated_at = ?
		WHERE id = ? `

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation one reservation by Id
func (m *mysqlDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `DELETE FROM reservations WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `UPDATE reservations SET processed = ? WHERE id = ?`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `SELECT id, room_name, created_at, updated_at FROM rooms ORDER BY room_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *mysqlDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `SELECT id, coalesce(reservation_id, 0) as res, restriction_id, room_id, start_date, end_date
			FROM room_restrictions WHERE ? <= end_date and ? >= start_date
			AND room_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationId,
			&r.RestrictionId,
			&r.RoomId,
			&r.StartDate,
			&r.EndDate,
		)

		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return restrictions, nil
}

// InsertBlockForRoom inserts a room restriction
func (m *mysqlDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?)
		`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// InsertBlockForRoom deletes a room restriction
func (m *mysqlDBRepo) DeleteBlockById(id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM room_restrictions WHERE id = ?;
		`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
