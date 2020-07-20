package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

var dbInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", Host, Port, User, Password, DBname, SSLmode)

func connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")

	return db, nil
}

func insertUnis(unis []*University) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStrings []string
	var valueArgs []interface{}
	for i, uni := range unis {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i * 9 + 1, i * 9 + 2, i * 9 + 3, i * 9 + 4, i * 9 + 5, i * 9 + 6, i * 9 + 7, i * 9 + 8, i * 9 + 9))
		valueArgs = append(valueArgs, uni.UniversityId)
		valueArgs = append(valueArgs, uni.Name)
		valueArgs = append(valueArgs, uni.Description)
		valueArgs = append(valueArgs, uni.Site)
		valueArgs = append(valueArgs, uni.Email)
		valueArgs = append(valueArgs, uni.Adress)
		valueArgs = append(valueArgs, uni.Phone)
		valueArgs = append(valueArgs, uni.MilitaryDep)
		valueArgs = append(valueArgs, uni.Dormitary)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO university VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
	}
}

func insertProfsNSpecs(profs map[Profile]bool, specsBach []*Speciality, specsSpec []*Speciality) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStringsProfs []string
	var valueArgsProfs []interface{}
	i := 0
	for p := range profs {
		valueStringsProfs = append(valueStringsProfs, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgsProfs = append(valueArgsProfs, p.ProfileId)
		valueArgsProfs = append(valueArgsProfs, p.Name)
		i++
	}

	sqlStmt := fmt.Sprintf("INSERT INTO profile VALUES %s;", strings.Join(valueStringsProfs, ","))
	if _, err = db.Exec(sqlStmt, valueArgsProfs...); err != nil {
		log.Println(err)
	}

	var valueStringsSpecs []string
	var valueArgsSpecs []interface{}
	i = 0
	for _, s := range specsBach {
		valueStringsSpecs = append(valueStringsSpecs, fmt.Sprintf("($%d, $%d, $%d, $%d)", i * 4 + 1, i * 4 + 2, i * 4 + 3, i * 4 + 4))
		valueArgsSpecs = append(valueArgsSpecs, s.SpecialityId)
		valueArgsSpecs = append(valueArgsSpecs, s.Name)
		valueArgsSpecs = append(valueArgsSpecs, s.Bachelor)
		valueArgsSpecs = append(valueArgsSpecs, s.ProfileId)
		i++
	}
	for _, s := range specsSpec {
		valueStringsSpecs = append(valueStringsSpecs, fmt.Sprintf("($%d, $%d, $%d, $%d)", i * 4 + 1, i * 4 + 2, i * 4 + 3, i * 4 + 4))
		valueArgsSpecs = append(valueArgsSpecs, s.SpecialityId)
		valueArgsSpecs = append(valueArgsSpecs, s.Name)
		valueArgsSpecs = append(valueArgsSpecs, s.Bachelor)
		valueArgsSpecs = append(valueArgsSpecs, s.ProfileId)
		i++
	}

	sqlStmt2 := fmt.Sprintf("INSERT INTO speciality VALUES %s;", strings.Join(valueStringsSpecs, ","))
	if _, err = db.Exec(sqlStmt2, valueArgsSpecs...); err != nil {
		log.Println(err)
	}

}

func getUnisIdsFromDb() []*University {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT university_id FROM university;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var unis []*University
	for rows.Next() {
		var university_id int
		err := rows.Scan(&university_id)
		if err != nil {
			log.Fatal(err)
		}

		uni := &University{
			UniversityId: university_id,
		}
		unis = append(unis, uni)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return unis
}

func insertFacs(facs []*Faculty) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStrings []string
	var valueArgs []interface{}
	for i, fac := range facs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i * 8 + 1, i * 8 + 2, i * 8 + 3, i * 8 + 4, i * 8 + 5, i * 8 + 6, i * 8 + 7, i * 8 + 8))
		valueArgs = append(valueArgs, fac.FacultyId)
		valueArgs = append(valueArgs, fac.Name)
		valueArgs = append(valueArgs, fac.Description)
		valueArgs = append(valueArgs, fac.Site)
		valueArgs = append(valueArgs, fac.Email)
		valueArgs = append(valueArgs, fac.Adress)
		valueArgs = append(valueArgs, fac.Phone)
		valueArgs = append(valueArgs, fac.UniversityId)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO faculty VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
	}
}