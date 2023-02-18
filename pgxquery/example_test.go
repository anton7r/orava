package pgxquery_test

import "github.com/jackc/pgx/v4/pgxpool"

func ExampleSelectNamed() {
	type User struct {
		ID       string `db:"user_id"`
		FullName string
		Email    string
		Age      int
	}

	type Table struct {
		Name string
	}

	api, _ := getAPI()

	db, _ := pgxpool.Connect(ctx, "example-connection-url")

	var users []*User
	if err := api.SelectNamed(
		ctx, db, &users, `SELECT user_id, full_name, email, age FROM :name`, &Table{Name: "users"},
	); err != nil {
		// Handle query or rows processing error.
	}
	// users variable now contains data from all rows.
}

func ExampleGetNamed() {
	type User struct {
		ID       string `db:"user_id"`
		FullName string
		Email    string
		Age      int
	}

	api, _ := getAPI()

	db, _ := pgxpool.Connect(ctx, "example-connection-url")

	var user User
	if err := api.GetNamed(
		ctx, db, &user, `SELECT full_name, email, age FROM users WHERE user_id=:user_id`, &User{ID: "bob"},
	); err != nil {
		// Handle query or rows processing error.
	}
	// user variable now contains data from all rows.
}

func ExampleExecNamed() {
	type User struct {
		ID       string `db:"user_id"`
		FullName string
		Email    string
		Age      int
	}

	db, _ := pgxpool.Connect(ctx, "example-connection-url")

	user := &User{
		ID:       "billy",
		FullName: "Billy Bob",
		Email:    "billy@example.com",
		Age:      50,
	}

	api, _ := getAPI()

	if _, err := api.ExecNamed(
		ctx, db, `INSERT INTO users (full_name, email, age, user_id) VALUES (:full_name, :email, :age, :user_id)`, user,
	); err != nil {
		// Handle exec processing error.
	}
	// user has now been inserted into the users table

	// let us now delete it
	if _, err := api.ExecNamed(
		ctx, db, `DELETE FROM users WHERE user_id = :ID`, user,
	); err != nil {
		// Handle exec processing error.
	}
	// user has now been deleted from the database
}
