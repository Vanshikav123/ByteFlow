package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// user struct conaining user credemtials
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

/*
	if bcrypt.Cost(hashedPassword) < 14 {
	    newHash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	    // Store newHash in the database
	}
*/
func (m *UserModel) Insert(name, email, password string) error {
	/*bcrypt.GenerateFromPassword([]byte(password), 12):
	This function hashes the provided password using the bcrypt algorithm.
	[]byte(password) converts the password string into a byte slice because bcrypt operates on byte slices.
	12 is the cost factor (computational complexity), which determines how expensive the hashing operation is.
	Higher values make it slower but more secure.*/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	/*
				sqlErr, ok := err.(*mysql.MySQLError) → Extracts MySQL error details directly.
				If err is of type mysql.MySQLError, errors.As() stores it in mySQLError
				errors.As(err, &target)
		       The errors.As() function is used to check if an error is of a specific type, and if it is,
		       it extracts that error into a variable of the specific type.

	*/
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	/*password is a string.
	bcrypt.GenerateFromPassword converts the password string into a hash and returns it as a []byte.
	bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost).
	This method converts the string into a hash and returns it as a byte slice.*/
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	/*EXISTS: This is a special SQL operator used to test if a subquery returns any rows.
	If the subquery returns at least one row, the EXISTS operator evaluates to true; otherwise,
	it evaluates to false*/
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
