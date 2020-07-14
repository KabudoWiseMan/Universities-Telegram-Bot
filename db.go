package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

func connect() error {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", Host, Port, User, Password, DBname, SSLmode)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Successfully connected!")

	return nil
}
