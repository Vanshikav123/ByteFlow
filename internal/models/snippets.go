package models

import (
	"database/sql"
	"errors"
	"time"
)

/*
interface=> An interface in Go defines a set of methods that a type must implement.
Enables polymorphism
can be implemented by multiple types

method=> A method in Go is a function that operates on a specific type (struct)
tied to one type
*/

// snippet struct to store paramaters of snippets
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// database model
type SnippetModel struct {
	DB *sql.DB
}

// insert ,get and latest methods interact with database to store snippets of text
// This is a method of SnippetModel, meaning it operates on an instance of SnippetModel.
// m.DB.Exec(...) executes the SQL statement.
// result is of type sql.Result, which contains metadata about the executed query.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Exec is a method from Go’s database/sql package used to execute SQL statements that do not return rows.
	//It's used for INSERT, UPDATE, DELETE, and other statements that modify data.
	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}
	//  LastInsertId() retrieves the ID of the last inserted row.
	/*LastInsertId() retrieves the ID of the most recently inserted row in a table with an auto-incrementing primary key.*/
	/*Other Similar Functions=>
	RowsAffected()	Returns the number of rows affected by an INSERT, UPDATE, or DELETE.
	QueryRow()	Executes a query that returns a single row.
	Query()	Executes a query that returns multiple rows.*/
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The id (which is of type int64) is converted to int and returned.
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	/*This creates a new Snippet struct on the heap and stores its memory address in s.
	  s is a pointer to a Snippet (*Snippet).
	  Since Get returns *Snippet, using a pointer allows efficient memory handling (we avoid copying the entire struct).*/
	s := &Snippet{}
	/*Scan fills variables with values from the SQL query.
	  Why &? Because Scan needs pointers to modify s.ID, s.Title, etc.
	  Without &, Scan wouldn’t be able to update the struct fields.*/

	/*Similar Functions to Scan
		  Scan(dest ...interface{})	Reads columns from QueryRow result into variables.
	Rows.Scan(dest ...interface{})	Reads multiple rows using Query().
	Row.Next()	Moves to the next row in a multi-row result.
	Row.Err()	Checks for errors in row iteration.
	*/
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	/*

	   In Latest(), s := &Snippet{} is created inside a loop, and we append multiple such pointers to a slice.
	   Each iteration creates a new Snippet and its pointer is added to the snippets slice.



	*/
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
