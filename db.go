package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
)

var dbInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", Host, Port, User, Password, DBname, SSLmode)

func connectToDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDb(db io.Closer) {
	if err := db.Close(); err != nil {
		log.Println("error closing db connection, error:", err)
	} else {
		log.Println("Data base connection closed")
	}
}

func updateUnisInDb(db *sql.DB, unis []*University) error {
	unisTmpTableQuery := "CREATE TEMPORARY TABLE temp_university (" +
		"university_id INT PRIMARY KEY, " +
		"name VARCHAR(300) NOT NULL, " +
		"description TEXT NOT NULL, " +
		"site VARCHAR(200) NOT NULL, " +
		"email VARCHAR(254) NOT NULL, " +
		"address VARCHAR(400) NOT NULL, " +
		"city_id SMALLINT NULL, " +
		"phone VARCHAR(300) NOT NULL, " +
		"military_dep BOOLEAN NOT NULL, " +
		"dormitary BOOLEAN NOT NULL" +
		");"
	if _, err := db.Exec(unisTmpTableQuery); err != nil {
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	for i, uni := range unis {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i * 10 + 1, i * 10 + 2, i * 10 + 3, i * 10 + 4, i * 10 + 5, i * 10 + 6, i * 10 + 7, i * 10 + 8, i * 10 + 9, i * 10 + 10))
		valueArgs = append(valueArgs, uni.UniversityId)
		valueArgs = append(valueArgs, uni.Name)
		valueArgs = append(valueArgs, uni.Description)
		valueArgs = append(valueArgs, uni.Site)
		valueArgs = append(valueArgs, uni.Email)
		valueArgs = append(valueArgs, uni.Address)
		valueArgs = append(valueArgs, uni.CityId)
		valueArgs = append(valueArgs, uni.Phone)
		valueArgs = append(valueArgs, uni.MilitaryDep)
		valueArgs = append(valueArgs, uni.Dormitary)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_university VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateUnisQuery := "INSERT INTO university " +
		"SELECT * FROM temp_university " +
		"ON CONFLICT (university_id) DO UPDATE " +
		"SET name = EXCLUDED.name, " +
		"description = EXCLUDED.description, " +
		"site = EXCLUDED.site, " +
		"email = EXCLUDED.email, " +
		"address = EXCLUDED.address, " +
		"city_id = EXCLUDED.city_id, " +
		"phone = EXCLUDED.phone, " +
		"military_dep = EXCLUDED.military_dep, " +
		"dormitary = EXCLUDED.dormitary;"
	if _, err := tx.Exec(updateUnisQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteUnisQuery := "DELETE FROM university WHERE university_id NOT IN (SELECT university_id FROM temp_university);"
	if _, err := tx.Exec(deleteUnisQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func updateProfsNSpecsInDb(db *sql.DB, profs []*Profile, specs []*Speciality) error {
	profsTmpTableQuery := "CREATE TEMPORARY TABLE temp_profile (" +
		"profile_id INT PRIMARY KEY, " +
		"name VARCHAR(200) NOT NULL" +
		");"
	if _, err := db.Exec(profsTmpTableQuery); err != nil {
		return err
	}

	specsTmpTableQuery := "CREATE TEMPORARY TABLE temp_speciality (" +
		"speciality_id INT PRIMARY KEY, " +
		"name VARCHAR(200) NOT NULL, " +
		"bachelor BOOLEAN NOT NULL, " +
		"profile_id INT NOT NULL" +
		");"
	if _, err := db.Exec(specsTmpTableQuery); err != nil {
		return err
	}

	var valueStringsProfs []string
	var valueArgsProfs []interface{}
	for i, p := range profs {
		valueStringsProfs = append(valueStringsProfs, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgsProfs = append(valueArgsProfs, p.ProfileId)
		valueArgsProfs = append(valueArgsProfs, p.Name)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_profile VALUES %s;", strings.Join(valueStringsProfs, ","))
	if _, err := db.Exec(sqlStmt, valueArgsProfs...); err != nil {
		return err
	}

	var valueStringsSpecs []string
	var valueArgsSpecs []interface{}
	for i, s := range specs {
		valueStringsSpecs = append(valueStringsSpecs, fmt.Sprintf("($%d, $%d, $%d, $%d)", i * 4 + 1, i * 4 + 2, i * 4 + 3, i * 4 + 4))
		valueArgsSpecs = append(valueArgsSpecs, s.SpecialityId)
		valueArgsSpecs = append(valueArgsSpecs, s.Name)
		valueArgsSpecs = append(valueArgsSpecs, s.Bachelor)
		valueArgsSpecs = append(valueArgsSpecs, s.ProfileId)
		i++
	}

	sqlStmt2 := fmt.Sprintf("INSERT INTO temp_speciality VALUES %s;", strings.Join(valueStringsSpecs, ","))
	if _, err := db.Exec(sqlStmt2, valueArgsSpecs...); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateProfsQuery := "INSERT INTO profile " +
		"SELECT * FROM temp_profile " +
		"ON CONFLICT (profile_id) DO UPDATE " +
		"SET name = EXCLUDED.name;"
	if _, err := tx.Exec(updateProfsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteProfsQuery := "DELETE FROM profile WHERE profile_id NOT IN (SELECT profile_id FROM temp_profile);"
	if _, err := tx.Exec(deleteProfsQuery); err != nil {
		tx.Rollback()
		return err
	}

	updateSpecsQuery := "INSERT INTO speciality " +
		"SELECT * FROM temp_speciality " +
		"ON CONFLICT (speciality_id) DO UPDATE " +
		"SET name = EXCLUDED.name, " +
		"bachelor = EXCLUDED.bachelor, " +
		"profile_id = EXCLUDED.profile_id;"
	if _, err := tx.Exec(updateSpecsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteSpecsQuery := "DELETE FROM speciality WHERE speciality_id NOT IN (SELECT speciality_id FROM temp_speciality);"
	if _, err := tx.Exec(deleteSpecsQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getUnisIdsNamesFromDb(db *sql.DB, withNames bool) ([]*UniversityInfo, error) {
	var unis []*UniversityInfo
	if withNames {
		rows, err := db.Query("SELECT university_id, name FROM university;")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var university_id int
			var name string
			err := rows.Scan(&university_id, &name)
			if err != nil {
				return nil, err
			}

			uni := &UniversityInfo{
				UniversityId: university_id,
				Name: name,
			}
			unis = append(unis, uni)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	} else {
		rows, err := db.Query("SELECT university_id FROM university;")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var university_id int
			err := rows.Scan(&university_id)
			if err != nil {
				return nil, err
			}

			uni := &UniversityInfo{
				UniversityId: university_id,
			}
			unis = append(unis, uni)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	return unis, nil
}

func updateFacsInDb(db *sql.DB, facs []*Faculty) error {
	facsTmpTableQuery := "CREATE TEMPORARY TABLE temp_faculty (" +
		"faculty_id INT PRIMARY KEY, " +
		"name VARCHAR(300) NOT NULL, " +
		"description TEXT NOT NULL, " +
		"site VARCHAR(200) NOT NULL, " +
		"email VARCHAR(254) NOT NULL, " +
		"address VARCHAR(400) NOT NULL, " +
		"phone VARCHAR(300) NOT NULL, " +
		"university_id INT NOT NULL" +
		");"
	if _, err := db.Exec(facsTmpTableQuery); err != nil {
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	for i, fac := range facs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i * 8 + 1, i * 8 + 2, i * 8 + 3, i * 8 + 4, i * 8 + 5, i * 8 + 6, i * 8 + 7, i * 8 + 8))
		valueArgs = append(valueArgs, fac.FacultyId)
		valueArgs = append(valueArgs, fac.Name)
		valueArgs = append(valueArgs, fac.Description)
		valueArgs = append(valueArgs, fac.Site)
		valueArgs = append(valueArgs, fac.Email)
		valueArgs = append(valueArgs, fac.Address)
		valueArgs = append(valueArgs, fac.Phone)
		valueArgs = append(valueArgs, fac.UniversityId)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_faculty VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateFacsQuery := "INSERT INTO faculty " +
		"SELECT * FROM temp_faculty " +
		"ON CONFLICT (faculty_id) DO UPDATE " +
		"SET name = EXCLUDED.name, " +
		"description = EXCLUDED.description, " +
		"site = EXCLUDED.site, " +
		"email = EXCLUDED.email, " +
		"address = EXCLUDED.address, " +
		"university_id = EXCLUDED.university_id;"
	if _, err := tx.Exec(updateFacsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteFacsQuery := "DELETE FROM faculty WHERE faculty_id NOT IN (SELECT faculty_id FROM temp_faculty);"
	if _, err := tx.Exec(deleteFacsQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getFacsIdsFromDb(db *sql.DB) ([]*Faculty, error) {
	rows, err := db.Query("SELECT faculty_id FROM faculty;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facs []*Faculty
	for rows.Next() {
		var faculty_id int
		err := rows.Scan(&faculty_id)
		if err != nil {
			return nil, err
		}

		fac := &Faculty{
			FacultyId: faculty_id,
		}
		facs = append(facs, fac)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return facs, nil
}

func updateSubjsInDb(db *sql.DB, subjs map[string]int) error {
	subjsTmpTableQuery := "CREATE TEMPORARY TABLE temp_subject (" +
		"subject_id SMALLINT PRIMARY KEY, " +
		"name VARCHAR(100) NOT NULL" +
		");"
	if _, err := db.Exec(subjsTmpTableQuery); err != nil {
		return err
	}

	var valueStringsSubjs []string
	var valueArgsSubjs []interface{}
	i := 0
	for s, k := range subjs {
		valueStringsSubjs = append(valueStringsSubjs, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgsSubjs = append(valueArgsSubjs, k)
		valueArgsSubjs = append(valueArgsSubjs, s)
		i++
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_subject VALUES %s;", strings.Join(valueStringsSubjs, ","))
	if _, err := db.Exec(sqlStmt, valueArgsSubjs...); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateSubjsQuery := "INSERT INTO subject " +
		"SELECT * FROM temp_subject " +
		"ON CONFLICT (subject_id) DO UPDATE " +
		"SET name = EXCLUDED.name;"
	if _, err := tx.Exec(updateSubjsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteSubjsQuery := "DELETE FROM subject WHERE subject_id NOT IN (SELECT subject_id FROM temp_subject);"
	if _, err := tx.Exec(deleteSubjsQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getRevSubjsMapFromDb(db *sql.DB) (map[string]int, error) {
	subjsRows, err := db.Query("SELECT * FROM subject;")
	if err != nil {
		return nil, err
	}
	defer subjsRows.Close()

	subjs := make(map[string]int)
	for subjsRows.Next() {
		var subject_id int
		var name string
		err := subjsRows.Scan(&subject_id, &name)
		if err != nil {
			return nil, err
		}

		subjs[name] = subject_id
	}
	err = subjsRows.Err()
	if err != nil {
		return nil, err
	}

	return subjs, nil
}

func insertTempProgs(db *sql.DB, progs []*Program) error {
	var valueStrings []string
	var valueArgs []interface{}
	for i, prog := range progs {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i * 15 + 1, i * 15 + 2, i * 15 + 3, i * 15 + 4, i * 15 + 5, i * 15 + 6, i * 15 + 7, i * 15 + 8, i * 15 + 9, i * 15 + 10, i * 15 + 11, i * 15 + 12, i * 15 + 13, i * 15 + 14, i * 15 + 15))
		valueArgs = append(valueArgs, prog.ProgramId)
		valueArgs = append(valueArgs, prog.ProgramNum)
		valueArgs = append(valueArgs, prog.Name)
		valueArgs = append(valueArgs, prog.Description)
		valueArgs = append(valueArgs, prog.FreePlaces)
		valueArgs = append(valueArgs, prog.PaidPlaces)
		valueArgs = append(valueArgs, prog.Fee)
		valueArgs = append(valueArgs, prog.FreePassPoints)
		valueArgs = append(valueArgs, prog.PaidPassPoints)
		valueArgs = append(valueArgs, prog.StudyForm)
		valueArgs = append(valueArgs, prog.StudyLanguage)
		valueArgs = append(valueArgs, prog.StudyBase)
		valueArgs = append(valueArgs, prog.StudyYears)
		valueArgs = append(valueArgs, prog.FacultyId)
		valueArgs = append(valueArgs, prog.SpecialityId)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_program VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	return nil
}

func insertTempMinPoints(db *sql.DB, minEgePoints []*MinEgePoints) error {
	var valueStrings []string
	var valueArgs []interface{}
	for i, minPoints := range minEgePoints {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, minPoints.ProgramId)
		valueArgs = append(valueArgs, minPoints.SubjectId)
		valueArgs = append(valueArgs, minPoints.MinPoints)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_min_ege_points VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	return nil
}

func insertTempEntrTests(db *sql.DB, entrTests []*EntranceTest) error {
	var valueStrings []string
	var valueArgs []interface{}
	for i, entrTest := range entrTests {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, entrTest.ProgramId)
		valueArgs = append(valueArgs, entrTest.TestName)
		valueArgs = append(valueArgs, entrTest.MinPoints)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_entrance_test VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	return nil
}

func createTempTablesForProgsNInfo(db *sql.DB) error {
	progsTmpTableQuery := "CREATE TEMPORARY TABLE temp_program (" +
		"program_id uuid PRIMARY KEY, " +
		"program_num INT NOT NULL, " +
		"name VARCHAR(400) NOT NULL, " +
		"description TEXT NOT NULL, " +
		"free_places SMALLINT NOT NULL, " +
		"paid_places SMALLINT NOT NULL, " +
		"fee MONEY NOT NULL, " +
		"free_pass_points SMALLINT NOT NULL, " +
		"paid_pass_points SMALLINT NOT NULL, " +
		"study_form VARCHAR(200) NOT NULL, " +
		"study_language VARCHAR(200) NOT NULL, " +
		"study_base VARCHAR(200) NOT NULL, " +
		"study_years VARCHAR(250) NOT NULL, " +
		"faculty_id INT NOT NULL, " +
		"speciality_id INT NOT NULL, " +
		"UNIQUE (faculty_id, program_num)" +
		");"
	if _, err := db.Exec(progsTmpTableQuery); err != nil {
		return err
	}

	minPointsTmpTableQuery := "CREATE TEMPORARY TABLE temp_min_ege_points (" +
		"program_id uuid NOT NULL, " +
		"subject_id SMALLINT NOT NULL, " +
		"min_points SMALLINT NOT NULL, " +
		"PRIMARY KEY (program_id, subject_id), " +
		"FOREIGN KEY (program_id) REFERENCES temp_program(program_id) ON UPDATE CASCADE" +
		");"
	if _, err := db.Exec(minPointsTmpTableQuery); err != nil {
		return err
	}

	entrTestsTmpTableQuery := "CREATE TEMPORARY TABLE temp_entrance_test (" +
		"program_id uuid NOT NULL, " +
		"test_name VARCHAR(200) NOT NULL, " +
		"min_points SMALLINT NOT NULL, " +
		"PRIMARY KEY (program_id, test_name), " +
		"FOREIGN KEY (program_id) REFERENCES temp_program(program_id) ON UPDATE CASCADE" +
		");"
	if _, err := db.Exec(entrTestsTmpTableQuery); err != nil {
		return err
	}

	return nil
}

func insertProgsNInfoToTempTables(db *sql.DB, progs []*Program, minEgePoints []*MinEgePoints, entrTests []*EntranceTest) error {
	if err := insertTempProgs(db, progs); err != nil {
		return err
	}
	if err := insertTempMinPoints(db, minEgePoints); err != nil {
		return err
	}
	if err := insertTempEntrTests(db, entrTests); err != nil {
		return err
	}

	return nil
}

func updateProgsNInfoInDb(db *sql.DB) error {
	matchProgramIdsQuery := "UPDATE temp_program " +
		"SET program_id = p.program_id " +
		"FROM program p " +
		"WHERE temp_program.program_num = p.program_num AND temp_program.faculty_id = p.faculty_id;"
	if _, err := db.Exec(matchProgramIdsQuery); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateProgsQuery := "INSERT INTO program " +
		"SELECT * FROM temp_program " +
		"ON CONFLICT (program_id) DO UPDATE " +
		"SET program_num = EXCLUDED.program_num, " +
		"name = EXCLUDED.name, " +
		"description = EXCLUDED.description, " +
		"free_places = EXCLUDED.free_places, " +
		"paid_places = EXCLUDED.paid_places, " +
		"fee = EXCLUDED.fee, " +
		"free_pass_points = EXCLUDED.free_pass_points, " +
		"paid_pass_points = EXCLUDED.paid_pass_points, " +
		"study_form = EXCLUDED.study_form, " +
		"study_language = EXCLUDED.study_language, " +
		"study_base = EXCLUDED.study_base, " +
		"study_years = EXCLUDED.study_years, " +
		"faculty_id = EXCLUDED.faculty_id, " +
		"speciality_id = EXCLUDED.speciality_id;"
	if _, err := tx.Exec(updateProgsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteProgsQuery := "DELETE FROM program WHERE program_id NOT IN (SELECT program_id FROM temp_program);"
	if _, err := tx.Exec(deleteProgsQuery); err != nil {
		tx.Rollback()
		return err
	}

	updateMinPointsQuery := "INSERT INTO min_ege_points " +
		"SELECT * FROM temp_min_ege_points " +
		"ON CONFLICT (program_id, subject_id) DO UPDATE " +
		"SET min_points = EXCLUDED.min_points;"
	if _, err := tx.Exec(updateMinPointsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteMinPointsQuery := "DELETE FROM min_ege_points WHERE (program_id, subject_id) NOT IN (SELECT program_id, subject_id FROM temp_min_ege_points);"
	if _, err := tx.Exec(deleteMinPointsQuery); err != nil {
		tx.Rollback()
		return err
	}

	updateEntrTestsQuery := "INSERT INTO entrance_test " +
		"SELECT * FROM temp_entrance_test " +
		"ON CONFLICT (program_id, test_name) DO UPDATE " +
		"SET min_points = EXCLUDED.min_points;"
	if _, err := tx.Exec(updateEntrTestsQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteEntrTestsQuery := "DELETE FROM entrance_test WHERE (program_id, test_name) NOT IN (SELECT program_id, test_name FROM temp_entrance_test);"
	if _, err := tx.Exec(deleteEntrTestsQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getSpecsIdsFromDb(db *sql.DB) ([]*Speciality, error) {
	rows, err := db.Query("SELECT speciality_id FROM speciality;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specs []*Speciality
	for rows.Next() {
		var speciality_id int
		err := rows.Scan(&speciality_id)
		if err != nil {
			return nil, err
		}

		spec := &Speciality{
			SpecialityId: speciality_id,
		}
		specs = append(specs, spec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func getUnisWithSitesFromDb(db *sql.DB) ([]*University, error) {
	rows, err := db.Query("SELECT university_id, site FROM university;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unis []*University
	for rows.Next() {
		var university_id int
		var site string
		err := rows.Scan(&university_id, &site)
		if err != nil {
			return nil, err
		}

		uni := &University{
			UniversityId: university_id,
			Site: site,
		}

		unis = append(unis, uni)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return unis, nil
}

func updateRatingQSInDb(db *sql.DB, ratingQS []*RatingQS) error {
	ratingQSTmpTableQuery := "CREATE TEMPORARY TABLE temp_rating_qs (" +
		"university_id INT PRIMARY KEY, " +
		"high_mark SMALLINT NOT NULL, " +
		"low_mark SMALLINT NOT NULL" +
		");"
	if _, err := db.Exec(ratingQSTmpTableQuery); err != nil {
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	for i, uniRatingQs := range ratingQS {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, uniRatingQs.UniversityId)
		valueArgs = append(valueArgs, uniRatingQs.HighMark)
		valueArgs = append(valueArgs, uniRatingQs.LowMark)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_rating_qs VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println("fuck")
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateRatingQSQuery := "INSERT INTO rating_qs " +
		"SELECT * FROM temp_rating_qs " +
		"ON CONFLICT (university_id) DO UPDATE " +
		"SET high_mark = EXCLUDED.high_mark, " +
		"low_mark = EXCLUDED.low_mark;"
	if _, err := tx.Exec(updateRatingQSQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteRatingQSQuery := "DELETE FROM rating_qs WHERE university_id NOT IN (SELECT university_id FROM temp_rating_qs);"
	if _, err := tx.Exec(deleteRatingQSQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func updateCitiesInDb(db *sql.DB, cities map[int]string) error {
	citiesTmpTableQuery := "CREATE TEMPORARY TABLE temp_city (" +
		"city_id SMALLINT PRIMARY KEY, " +
		"name VARCHAR(100)" +
		");"
	if _, err := db.Exec(citiesTmpTableQuery); err != nil {
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	i := 0
	for k, city := range cities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgs = append(valueArgs, k)
		valueArgs = append(valueArgs, city)
		i++
	}

	sqlStmt := fmt.Sprintf("INSERT INTO temp_city VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := db.Exec(sqlStmt, valueArgs...); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	updateCitiesQuery := "INSERT INTO city " +
		"SELECT * FROM temp_city " +
		"ON CONFLICT (city_id) DO UPDATE " +
		"SET name = EXCLUDED.name;"
	if _, err := tx.Exec(updateCitiesQuery); err != nil {
		tx.Rollback()
		return err
	}

	deleteCitiesQuery := "DELETE FROM city WHERE city_id NOT IN (SELECT city_id FROM temp_city);"
	if _, err := tx.Exec(deleteCitiesQuery); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func getCountFromDb(db *sql.DB, from string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM " + from + ";").Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func getUnisQSNumFromDb(db *sql.DB) (int, error) {
	return getCountFromDb(db, "rating_qs")
}

func makeQSMark(high_mark int, low_mark int) string {
	var mark string
	if high_mark == low_mark {
		mark = strconv.Itoa(high_mark)
	} else {
		mark = strconv.Itoa(high_mark) + "-" + strconv.Itoa(low_mark)
	}

	return mark
}

func getUnisQSPageFromDb(db *sql.DB, offset string) ([]*UniversityQS, error) {
	rows, err := db.Query("SELECT u.university_id, u.name, rq.high_mark, rq.low_mark FROM university u JOIN rating_qs rq on u.university_id = rq.university_id ORDER BY high_mark LIMIT 5 OFFSET " + offset + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var universitiesQS []*UniversityQS
	for rows.Next() {
		var university_id, high_mark, low_mark int
		var name string
		err := rows.Scan(&university_id, &name, &high_mark, &low_mark)
		if err != nil {
			return nil, err
		}

		mark := makeQSMark(high_mark, low_mark)

		universityQs := &UniversityQS{
			UniversityId: university_id,
			Name: name,
			Mark: mark,
		}

		universitiesQS = append(universitiesQS, universityQs)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return universitiesQS, nil
}

func getUniFromDb(db *sql.DB, uniId string) (*UniversityInfo, error) {
	uni := &UniversityInfo{}
	query := "SELECT u.university_id, u.name, u.description, u.site, u.email, COALESCE('г. ' || c.name || ', ' || u.address, u.address), u.phone, u.military_dep, u.dormitary FROM university u " +
		"LEFT JOIN city c on c.city_id = u.city_id " +
		"WHERE university_id = " + uniId + ";"
	err := db.QueryRow(query).Scan(&uni.UniversityId, &uni.Name, &uni.Description, &uni.Site, &uni.Email, &uni.Address, &uni.Phone, &uni.MilitaryDep, &uni.Dormitary)
	if err != nil {
		return nil, err
	}

	return uni, nil
}

func getUniQSRateFromDb(db *sql.DB, uniId string) (string, error) {
	var high_mark, low_mark int
	err := db.QueryRow("SELECT high_mark, low_mark FROM rating_qs WHERE university_id = " + uniId + ";").Scan(&high_mark, &low_mark)
	if err != nil {
		if err = db.Ping(); err != nil {
			return "", err
		}

		return "", nil
	}

	mark := makeQSMark(high_mark, low_mark)

	return mark, nil
}

func getFacsNumFromDb(db *sql.DB, uniId string) (int, error) {
	return getCountFromDb(db, "faculty WHERE university_id = " + uniId)
}

func getFacsPageFromDb(db *sql.DB, uniId string, offset string) ([]*Faculty, error) {
	rows, err := db.Query("SELECT faculty_id, name FROM faculty WHERE university_id = " + uniId + " LIMIT 5 OFFSET " + offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facs []*Faculty
	for rows.Next() {
		var faculty_id int
		var name string
		err := rows.Scan(&faculty_id, &name)
		if err != nil {
			return nil, err
		}

		uniIdNum, _ := strconv.Atoi(uniId)
		fac := &Faculty{
			UniversityId: uniIdNum,
			FacultyId: faculty_id,
			Name: name,
		}

		facs = append(facs, fac)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return facs, nil
}

func getFacFromDb(db *sql.DB, facId string) (*Faculty, error) {
	fac := &Faculty{}
	err := db.QueryRow("SELECT * FROM faculty WHERE faculty_id = " + facId + ";").Scan(&fac.FacultyId, &fac.Name, &fac.Description, &fac.Site, &fac.Email, &fac.Address, &fac.Phone, &fac.UniversityId)
	if err != nil {
		return nil, err
	}

	return fac, nil
}

func getFindUnisNumFromDb(db *sql.DB, query string) (int, error) {
	return getCountFromDb(db, "university_name_descr_vector WHERE name_descr_vector @@ plainto_tsquery('russian', '" + query + "')")
}

func getUnisIdsNNamesFromDb(db *sql.DB, query string) ([]*UniversityInfo, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unis []*UniversityInfo
	for rows.Next() {
		var university_id int
		var name string
		err := rows.Scan(&university_id, &name)
		if err != nil {
			return nil, err
		}

		uni := &UniversityInfo{
			UniversityId: university_id,
			Name: name,
		}

		unis = append(unis, uni)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return unis, nil
}

func findUnisInDb(db *sql.DB, query string, offset string) ([]*UniversityInfo, error) {
	dbQuery := "SELECT u.university_id, u.name FROM university u " +
		"JOIN (" +
		"SELECT university_id, ts_rank(name_descr_vector, plainto_tsquery('russian', '" + query + "')) " +
		"FROM university_name_descr_vector " +
		"WHERE name_descr_vector @@ plainto_tsquery('russian', '" + query + "') " +
		"ORDER BY ts_rank(name_descr_vector, plainto_tsquery('russian', '" + query + "')) DESC " +
		"LIMIT 5 OFFSET " + offset +
		") l ON (u.university_id = l.university_id);"

	return getUnisIdsNNamesFromDb(db, dbQuery)
}

func getUniProfsNumFromDb(db *sql.DB, uniId string) (int, error) {
	from := "(" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		") l"
	return getCountFromDb(db, from)
}

func getProfsFromDb(db *sql.DB, query string) ([]*Profile, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profs []*Profile
	for rows.Next() {
		var profile_id int
		var name string
		err := rows.Scan(&profile_id, &name)
		if err != nil {
			return nil, err
		}

		prof := &Profile{
			ProfileId: profile_id,
			Name: name,
		}

		profs = append(profs, prof)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return profs, nil
}

func getUniProfsPageFromDb(db *sql.DB, uniId string, offset string) ([]*Profile, error) {
	query := "SELECT p.* FROM profile p " +
		"JOIN (" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		") l ON (p.profile_id = l.profile_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getProfsFromDb(db, query)
}

func getFacProfsNumFromDb(db *sql.DB, facId string) (int, error) {
	from := "(" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		") l"
	return getCountFromDb(db, from)
}

func getFacProfsPageFromDb(db *sql.DB, facId string, offset string) ([]*Profile, error) {
	query := "SELECT p.* FROM profile p " +
		"JOIN (" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		") l ON (p.profile_id = l.profile_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getProfsFromDb(db, query)
}

func getProfFromDb(db *sql.DB, profId string) (*Profile, error) {
	prof := &Profile{}
	err := db.QueryRow("SELECT * FROM profile WHERE profile_id = " + profId + ";").Scan(&prof.ProfileId, &prof.Name)
	if err != nil {
		return nil, err
	}

	return prof, nil
}

func getSpecsFromDb(db *sql.DB, query string) ([]*Speciality, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specs []*Speciality
	for rows.Next() {
		var speciality_id, profile_id int
		var name string
		var bachelor bool
		err := rows.Scan(&speciality_id, &name, &bachelor, &profile_id)
		if err != nil {
			return nil, err
		}

		spec := &Speciality{
			SpecialityId: speciality_id,
			Name: name,
			Bachelor: bachelor,
			ProfileId: profile_id,
		}

		specs = append(specs, spec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return specs, nil
}

func getUniSpecsNumFromDb(db *sql.DB, uniId string, profId string) (int, error) {
	from := "(" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND s.profile_id = " + profId +
		") l"
	return getCountFromDb(db, from)
}

func getUniSpecsPageFromDb(db *sql.DB, uniId string, profId string, offset string) ([]*Speciality, error) {
	query := "SELECT s.* FROM speciality s " +
		"JOIN (" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND s.profile_id = " + profId +
		") l ON (s.speciality_id = l.speciality_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getSpecsFromDb(db, query)
}

func getFacSpecsNumFromDb(db *sql.DB, facId string, profId string) (int, error) {
	from := "(" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND s.profile_id = " + profId +
		") l"
	return getCountFromDb(db, from)
}

func getFacSpecsPageFromDb(db *sql.DB, facId string, profId string, offset string) ([]*Speciality, error) {
	query := "SELECT s.* FROM speciality s " +
		"JOIN (" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND s.profile_id = " + profId +
		") l ON (s.speciality_id = l.speciality_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getSpecsFromDb(db, query)
}

func getSpecFromDb(db *sql.DB, specId string) (*Speciality, error) {
	spec := &Speciality{}
	err := db.QueryRow("SELECT * FROM speciality WHERE speciality_id = " + specId + ";").Scan(&spec.SpecialityId, &spec.Name, &spec.Bachelor, &spec.ProfileId)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func getProgsInfoFromDb(db *sql.DB, query string) ([]*ProgramPreview, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progs []*ProgramPreview
	for rows.Next() {
		var program_id uuid.UUID
		var speciality_id int
		var name, spec_name string
		err := rows.Scan(&program_id, &name, &speciality_id, &spec_name)
		if err != nil {
			return nil, err
		}

		prog := &ProgramPreview{
			ProgramId: program_id,
			Name: name,
			SpecialityId: speciality_id,
			SpecialityName: spec_name,
		}

		progs = append(progs, prog)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return progs, nil
}

func getUniProgsNumFromDb(db *sql.DB, uniId string) (int, error) {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId
	return getCountFromDb(db, from)
}

func getUniProgsPageFromDb(db *sql.DB, uniId string, offset string) ([]*ProgramPreview, error) {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id, s.name FROM program pr " +
		"JOIN speciality s ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(db, query)
}

func getFacProgsNumFromDb(db *sql.DB, facId string) (int, error) {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId
	return getCountFromDb(db, from)
}

func getFacProgsPageFromDb(db *sql.DB, facId string, offset string) ([]*ProgramPreview, error) {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id, s.name FROM program pr " +
		"JOIN speciality s ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(db, query)
}

func getUniSpecProgsNumFromDb(db *sql.DB, uniId string, specId string) (int, error) {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND pr.speciality_id = " + specId
	return getCountFromDb(db, from)
}

func getUniSpecProgsPageFromDb(db *sql.DB, uniId string, specId string, offset string) ([]*ProgramPreview, error) {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id, s.name FROM program pr " +
		"JOIN speciality s ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND pr.speciality_id = " + specId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(db, query)
}

func getFacSpecProgsNumFromDb(db *sql.DB, facId string, specId string) (int, error) {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND pr.speciality_id = " + specId
	return getCountFromDb(db, from)
}

func getFacSpecProgsPageFromDb(db *sql.DB, facId string, specId string, offset string) ([]*ProgramPreview, error) {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id, s.name FROM program pr " +
		"JOIN speciality s ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND pr.speciality_id = " + specId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(db, query)
}

func getProgInfoFromDb(db *sql.DB, progId string) (*ProgramInfo, error) {
	prog := &ProgramInfo{}
	query := "SELECT pr.program_id, pr.program_num, pr.name, pr.description, pr.free_places, pr.paid_places, pr.fee::numeric::int8, pr.free_pass_points, pr.paid_pass_points, pr.study_form, pr.study_language, pr.study_base, pr.study_years, pr.faculty_id, pr.speciality_id, s2.name, s2.bachelor, COALESCE(l.ege, '') as eges, COALESCE(l2.entrs, '') as entrs FROM program pr " +
		"LEFT JOIN (" +
		"SELECT m.program_id, string_agg(s.name || ' ' || m.min_points::text, E'\\n') as ege FROM min_ege_points m " +
		"JOIN subject s ON (m.subject_id = s.subject_id) " +
		"WHERE m.program_id = '" + progId + "' " +
		"GROUP BY m.program_id" +
		") l ON (pr.program_id = l.program_id) " +
		"LEFT JOIN (" +
		"SELECT et.program_id, string_agg(et.test_name || ' ' || et.min_points::text, E'\\n') as entrs FROM entrance_test et " +
		"WHERE et.program_id = '" + progId + "' " +
		"GROUP BY et.program_id" +
		") l2 ON (pr.program_id = l2.program_id) " +
		"LEFT JOIN speciality s2 ON (pr.speciality_id = s2.speciality_id) " +
		"WHERE pr.program_id = '" + progId + "';"
	err := db.QueryRow(query).Scan(&prog.ProgramId, &prog.ProgramNum, &prog.Name, &prog.Description, &prog.FreePlaces, &prog.PaidPlaces, &prog.Fee, &prog.FreePassPoints, &prog.PaidPassPoints, &prog.StudyForm, &prog.StudyLanguage, &prog.StudyBase, &prog.StudyYears, &prog.FacultyId, &prog.SpecialityId, &prog.SpecialityName, &prog.Bachelor, &prog.EGEs, &prog.EntranceTests)
	if err != nil {
		return nil, err
	}

	return prog, nil
}

func getUniOfFacFromDb(db *sql.DB, facId string) (*UniversityInfo, error) {
	uni := &UniversityInfo{}
	err := db.QueryRow("SELECT u.university_id, u.name FROM university u JOIN faculty f ON (u.university_id = f.university_id) WHERE faculty_id = " + facId + ";").Scan(&uni.UniversityId, &uni.Name)
	if err != nil {
		return nil, err
	}

	return uni, nil
}

func getCitiesNumFromDb(db *sql.DB) (int, error) {
	from := "city"
	return getCountFromDb(db, from)
}

func getAllCitiesFromDb(db *sql.DB) ([]*City, error) {
	rows, err := db.Query("SELECT * FROM city;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []*City
	for rows.Next() {
		var city_id int
		var name string
		err := rows.Scan(&city_id, &name)
		if err != nil {
			return nil, err
		}

		city := &City{
			CityId: city_id,
			Name: name,
		}

		cities = append(cities, city)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func getCitiesFromDb(db *sql.DB, offset string) ([]*City, error) {
	rows, err := db.Query("SELECT * FROM city ORDER BY name LIMIT 5 OFFSET " + offset + ";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []*City
	for rows.Next() {
		var city_id int
		var name string
		err := rows.Scan(&city_id, &name)
		if err != nil {
			return nil, err
		}

		city := &City{
			CityId: city_id,
			Name: name,
		}

		cities = append(cities, city)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func getCityNameFromDb(db *sql.DB, cityId int) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM city WHERE city_id = " + strconv.Itoa(cityId) + ";").Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func getProfsNumFromDb(db *sql.DB) (int, error) {
	from := "profile"
	return getCountFromDb(db, from)
}

func getProfsPageFromDb(db *sql.DB, offset string) ([]*Profile, error) {
	query := "SELECT * FROM profile ORDER BY name LIMIT 5 OFFSET " + offset
	return getProfsFromDb(db, query)
}

func getSpecsNumFromDb(db *sql.DB, profId string) (int, error) {
	from := "speciality WHERE profile_id = " + profId
	return getCountFromDb(db, from)
}

func getSpecsPageFromDb(db *sql.DB, offset string, profId string) ([]*Speciality, error) {
	query := "SELECT * FROM speciality WHERE profile_id = " + profId + " ORDER BY name LIMIT 5 OFFSET " + offset
	return getSpecsFromDb(db, query)
}

func getSubjsMapFromDb(db *sql.DB) (map[int]string, error) {
	subjsRows, err := db.Query("SELECT * FROM subject;")
	if err != nil {
		return nil, err
	}
	defer subjsRows.Close()

	subjs := make(map[int]string)
	for subjsRows.Next() {
		var subject_id int
		var name string
		err := subjsRows.Scan(&subject_id, &name)
		if err != nil {
			return nil, err
		}

		subjs[subject_id] = name
	}
	err = subjsRows.Err()
	if err != nil {
		return nil, err
	}

	return subjs, nil
}

func getSubjsNumFromDb(db *sql.DB, user *UserInfo) (int, error) {
	from := "subject"
	if len(user.Eges) != 0 {
		from += " WHERE "
	}
	for i, ege := range user.Eges {
		if i == len(user.Eges) - 1 {
			from += "subject_id != " + strconv.Itoa(ege.SubjId)
		} else {
			from += "subject_id != " + strconv.Itoa(ege.SubjId) + "AND "
		}
	}

	return getCountFromDb(db, from)
}

func getSubjsFromDb(db *sql.DB, offset string, user *UserInfo) ([]*Subject, error) {
	var except string
	if len(user.Eges) != 0 {
		except += "WHERE "
	}
	for i, ege := range user.Eges {
		if i == len(user.Eges) - 1 {
			except += "subject_id != " + strconv.Itoa(ege.SubjId)
		} else {
			except += "subject_id != " + strconv.Itoa(ege.SubjId) + "AND "
		}
	}

	subjsRows, err := db.Query("SELECT * FROM subject " + except + " ORDER BY subject_id LIMIT 5 OFFSET " + offset + ";")
	if err != nil {
		return nil, err
	}
	defer subjsRows.Close()

	var subjs []*Subject
	for subjsRows.Next() {
		var subject_id int
		var name string
		err := subjsRows.Scan(&subject_id, &name)
		if err != nil {
			return nil, err
		}

		subj := &Subject{
			SubjectId: subject_id,
			Name: name,
		}

		subjs = append(subjs, subj)
	}
	err = subjsRows.Err()
	if err != nil {
		return nil, err
	}

	return subjs, nil
}

func getSubjNameFromDb(db *sql.DB, subjId int) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM subject WHERE subject_id = " + strconv.Itoa(subjId) + ";").Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func makeSearchInnerQueryForDb(user *UserInfo) string {
	var conds []string
	from := "SELECT DISTINCT u.university_id FROM university u " +
		"JOIN faculty f ON (u.university_id = f.university_id) " +
		"JOIN program p ON (f.faculty_id = p.faculty_id) "

	var pointsSum uint64
	if len(user.Eges) != 0 {
		for i, ege := range user.Eges {
			pointsSum += ege.MinPoints
			iStr := strconv.Itoa(i)
			from += "JOIN min_ege_points m" + iStr + " ON p.program_id = m" + iStr + ".program_id "
			conds = append(conds, "m" + iStr + ".subject_id = " + strconv.Itoa(ege.SubjId))
			if ege.MinPoints != 100 {
				conds = append(conds, "m" + iStr + ".min_points <= " + strconv.Itoa(int(ege.MinPoints)))
			}
		}
		from += "JOIN (" +
			"SELECT program_id FROM min_ege_points " +
			"GROUP BY program_id " +
			"HAVING COUNT(program_id) = " + strconv.Itoa(len(user.Eges)) +
			") l2 ON p.program_id = l2.program_id "
	}

	if !user.EntryTest {
		conds = append(conds, "p.program_id NOT IN (SELECT DISTINCT program_id FROM entrance_test)")
	}

	var feeConds string
	if user.Fee == 0 {
		if pointsSum > 0 {
			feeConds = "p.free_pass_points <= " + strconv.Itoa(int(pointsSum)) + " AND p.free_places > 0"
		} else {
			feeConds = "p.free_places > 0"
		}
	} else if user.Fee == math.MaxUint64 {
		if pointsSum > 0 {
			feeConds = "(p.free_pass_points <= " + strconv.Itoa(int(pointsSum)) + " AND p.free_places > 0 OR " +
				"p.paid_pass_points <= " + strconv.Itoa(int(pointsSum)) + " AND p.paid_places > 0)"
		}
	} else {
		if pointsSum > 0 {
			feeConds = "(p.free_pass_points <= " + strconv.Itoa(int(pointsSum)) + " AND p.free_places > 0 OR " +
				"p.paid_pass_points <= " + strconv.Itoa(int(pointsSum)) + " AND p.paid_places > 0 AND p.fee <= " + strconv.Itoa(int(user.Fee)) + "::money)"
		} else {
			feeConds = "p.fee <= " + strconv.Itoa(int(user.Fee)) + "::money"
		}
	}

	if feeConds != "" {
		conds = append(conds, feeConds)
	}

	if user.ProfileId != 0 {
		if user.SpecialityId != 0 {
			conds = append(conds, "p.speciality_id = " + strconv.Itoa(user.SpecialityId))
		} else {
			from += "JOIN speciality s ON (p.speciality_id = s.speciality_id) "
			conds = append(conds, "s.profile_id = " + strconv.Itoa(user.ProfileId))
		}
	}

	if user.City != 0 {
		conds = append(conds, "u.city_id = " + strconv.Itoa(user.City))
	}

	if user.Dormatary {
		conds = append(conds, "u.dormitary")
	}

	if user.MilitaryDep {
		conds = append(conds, "u.military_dep")
	}

	var wholeCond string
	if len(conds) != 0 {
		wholeCond += "WHERE " + strings.Join(conds, " AND ")
	}

	return from + wholeCond
}

func getSearchUnisNumFromDb(db *sql.DB, from string) (int, error) {
	return getCountFromDb(db, "(" + from + ") l")
}

func searchUnisInDb(db *sql.DB, innerQuery string, offset string) ([]*UniversityInfo, error) {
	query := "SELECT u.university_id, u.name FROM university u " +
		"JOIN (" + innerQuery +
		") l ON (u.university_id = l.university_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	log.Println("QUERY:", query)
	return getUnisIdsNNamesFromDb(db, query)
}