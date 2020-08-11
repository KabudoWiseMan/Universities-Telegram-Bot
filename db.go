package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
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

func insertProfsNSpecs(profs []*Profile, specs []*Speciality) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStringsProfs []string
	var valueArgsProfs []interface{}
	for i, p := range profs {
		valueStringsProfs = append(valueStringsProfs, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgsProfs = append(valueArgsProfs, p.ProfileId)
		valueArgsProfs = append(valueArgsProfs, p.Name)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO profile VALUES %s;", strings.Join(valueStringsProfs, ","))
	if _, err = db.Exec(sqlStmt, valueArgsProfs...); err != nil {
		log.Println(err)
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

	sqlStmt2 := fmt.Sprintf("INSERT INTO speciality VALUES %s;", strings.Join(valueStringsSpecs, ","))
	if _, err = db.Exec(sqlStmt2, valueArgsSpecs...); err != nil {
		log.Println(err)
	}
}

func getUnisIdsNamesFromDb(withNames bool) []*University {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var unis []*University
	if withNames {
		rows, err := db.Query("SELECT university_id, name FROM university;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var university_id int
			var name string
			err := rows.Scan(&university_id, &name)
			if err != nil {
				log.Fatal(err)
			}

			uni := &University{
				UniversityId: university_id,
				Name: name,
			}
			unis = append(unis, uni)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows, err := db.Query("SELECT university_id FROM university;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

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

func getFacsIdsFromDb() []*Faculty {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT faculty_id FROM faculty;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var facs []*Faculty
	for rows.Next() {
		var faculty_id int
		err := rows.Scan(&faculty_id)
		if err != nil {
			log.Fatal(err)
		}

		fac := &Faculty{
			FacultyId: faculty_id,
		}
		facs = append(facs, fac)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return facs
}

//func insertStudyForms(studyForms map[string]int) {
//	db, err := sql.Open("postgres", dbInfo)
//	if err != nil {
//		log.Fatal("Couldn't connect to db")
//	}
//	defer db.Close()
//
//	var valueStringsStudyForms []string
//	var valueArgsStudyForms []interface{}
//	i := 0
//	for f, k := range studyForms {
//		valueStringsStudyForms = append(valueStringsStudyForms, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
//		valueArgsStudyForms = append(valueArgsStudyForms, k)
//		valueArgsStudyForms = append(valueArgsStudyForms, f)
//		i++
//	}
//
//	sqlStmt2 := fmt.Sprintf("INSERT INTO study_form VALUES %s;", strings.Join(valueStringsStudyForms, ","))
//	if _, err = db.Exec(sqlStmt2, valueArgsStudyForms...); err != nil {
//		log.Println(err)
//	}
//}

func insertSubjs(subjs map[string]int) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStringsSubjs []string
	var valueArgsSubjs []interface{}
	i := 0
	for s, k := range subjs {
		valueStringsSubjs = append(valueStringsSubjs, fmt.Sprintf("($%d, $%d)", i * 2 + 1, i * 2 + 2))
		valueArgsSubjs = append(valueArgsSubjs, k)
		valueArgsSubjs = append(valueArgsSubjs, s)
		i++
	}

	sqlStmt := fmt.Sprintf("INSERT INTO subject VALUES %s;", strings.Join(valueStringsSubjs, ","))
	if _, err = db.Exec(sqlStmt, valueArgsSubjs...); err != nil {
		log.Println(err)
	}
}

//func getStudyFormsFromDb() map[string]int {
//	db, err := sql.Open("postgres", dbInfo)
//	if err != nil {
//		log.Fatal("Couldn't connect to db")
//	}
//	defer db.Close()
//
//	studyFormsRows, err := db.Query("SELECT * FROM study_form;")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer studyFormsRows.Close()
//
//	studyForms := make(map[string]int)
//	for studyFormsRows.Next() {
//		var study_form_id int
//		var name string
//		err := studyFormsRows.Scan(&study_form_id, &name)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		studyForms[name] = study_form_id
//	}
//	err = studyFormsRows.Err()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return studyForms
//}

func getSubjsFromDb() map[string]int {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	subjsRows, err := db.Query("SELECT * FROM subject;")
	if err != nil {
		log.Fatal(err)
	}
	defer subjsRows.Close()

	subjs := make(map[string]int)
	for subjsRows.Next() {
		var subject_id int
		var name string
		err := subjsRows.Scan(&subject_id, &name)
		if err != nil {
			log.Fatal(err)
		}

		subjs[name] = subject_id
	}
	err = subjsRows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return subjs
}

func insertProgs(tx *sql.Tx, progs []*Program) {
	//db, err := sql.Open("postgres", dbInfo)
	//if err != nil {
	//	log.Fatal("Couldn't connect to db")
	//}
	//defer db.Close()

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

	sqlStmt := fmt.Sprintf("INSERT INTO program VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
		tx.Rollback()
	}
}

func insertMinPoints(tx *sql.Tx, minEgePoints []*MinEgePoints) {
	//db, err := sql.Open("postgres", dbInfo)
	//if err != nil {
	//	log.Fatal("Couldn't connect to db")
	//}
	//defer db.Close()

	var valueStrings []string
	var valueArgs []interface{}
	for i, minPoints := range minEgePoints {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, minPoints.ProgramId)
		valueArgs = append(valueArgs, minPoints.SubjectId)
		valueArgs = append(valueArgs, minPoints.MinPoints)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO min_ege_points VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
		tx.Rollback()
	}
}

func insertEntrTests(tx *sql.Tx, entrTests []*EntranceTest) {
	//db, err := sql.Open("postgres", dbInfo)
	//if err != nil {
	//	log.Fatal("Couldn't connect to db")
	//}
	//defer db.Close()

	var valueStrings []string
	var valueArgs []interface{}
	for i, entrTest := range entrTests {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, entrTest.ProgramId)
		valueArgs = append(valueArgs, entrTest.TestName)
		valueArgs = append(valueArgs, entrTest.MinPoints)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO entrance_test VALUES %s;", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
		tx.Rollback()
	}
}

func insertProgsNInfo(progs []*Program, minEgePoints []*MinEgePoints, entrTests []*EntranceTest) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Couldn't begin the transaction")
	}

	insertProgs(tx, progs)
	insertMinPoints(tx, minEgePoints)

	if len(entrTests) != 0 {
		insertEntrTests(tx, entrTests)
	}

	err = tx.Commit()
	if err != nil {
		log.Println("couldn't commit the transaction")
	}
}

func getSpecsIdsFromDb() []*Speciality {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT speciality_id FROM speciality;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var specs []*Speciality
	for rows.Next() {
		var speciality_id int
		err := rows.Scan(&speciality_id)
		if err != nil {
			log.Fatal(err)
		}

		spec := &Speciality{
			SpecialityId: speciality_id,
		}
		specs = append(specs, spec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return specs
}

func getUniIdFromDb(uniSite string) int {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT university_id FROM university WHERE site LIKE" + "'%www." + uniSite + ".ru%' LIMIT 1;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var university_id int
	for rows.Next() {
		err := rows.Scan(&university_id)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if university_id == 0 {
		rows2, err := db.Query("SELECT university_id FROM university WHERE site LIKE" + "'%" + uniSite + ".ru%' LIMIT 1;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows2.Close()

		for rows2.Next() {
			err := rows2.Scan(&university_id)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = rows2.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	return university_id
}

func insertRatingQS(ratingQS []*RatingQS) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var valueStrings []string
	var valueArgs []interface{}
	for i, uniRatingQs := range ratingQS {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i * 3 + 1, i * 3 + 2, i * 3 + 3))
		valueArgs = append(valueArgs, uniRatingQs.UniversityId)
		valueArgs = append(valueArgs, uniRatingQs.HighMark)
		valueArgs = append(valueArgs, uniRatingQs.LowMark)
	}

	sqlStmt := fmt.Sprintf("INSERT INTO rating_qs VALUES %s;", strings.Join(valueStrings, ","))
	if _, err = db.Exec(sqlStmt, valueArgs...); err != nil {
		log.Println(err)
	}
}

func getCountFromDb(from string) int {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + from + ";").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	return count
}

func getUnisQSNumFromDb() int {
	return getCountFromDb("rating_qs")
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

func getUnisQSPageFromDb(offset string) []*UniversityQS {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT u.university_id, u.name, rq.high_mark, rq.low_mark FROM university u JOIN rating_qs rq on u.university_id = rq.university_id ORDER BY high_mark LIMIT 5 OFFSET " + offset + ";")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var universitiesQS []*UniversityQS
	for rows.Next() {
		var university_id, high_mark, low_mark int
		var name string
		err := rows.Scan(&university_id, &name, &high_mark, &low_mark)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	return universitiesQS
}

func getUniFromDb(uniId string) University {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var university_id int
	var name, description, site, email, adress, phone string
	var military_dep, dormitary bool
	err = db.QueryRow("SELECT * FROM university WHERE university_id = " + uniId + ";").Scan(&university_id, &name, &description, &site, &email, &adress, &phone, &military_dep, &dormitary)
	if err != nil {
		log.Fatal(err)
	}

	uni := University{
		UniversityId: university_id,
		Name: name,
		Description: description,
		Site: site,
		Email: email,
		Adress: adress,
		Phone: phone,
		MilitaryDep: military_dep,
		Dormitary: dormitary,
	}

	return uni
}

func getUniQSRateFromDb(uniId string) string {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var high_mark, low_mark int
	err = db.QueryRow("SELECT high_mark, low_mark FROM rating_qs WHERE university_id = " + uniId + ";").Scan(&high_mark, &low_mark)
	if err != nil {
		return ""
	}

	mark := makeQSMark(high_mark, low_mark)

	return mark
}

func getFacsNumFromDb(uniId string) int {
	return getCountFromDb("faculty WHERE university_id = " + uniId)
}

func getFacsPageFromDb(uniId string, offset string) []*Faculty {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query("SELECT faculty_id, name FROM faculty WHERE university_id = " + uniId + " LIMIT 5 OFFSET " + offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var facs []*Faculty
	for rows.Next() {
		var faculty_id int
		var name string
		err := rows.Scan(&faculty_id, &name)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	return facs
}

func getFacFromDb(facId string) Faculty {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var university_id, faculty_id int
	var name, description, site, email, adress, phone string
	err = db.QueryRow("SELECT * FROM faculty WHERE faculty_id = " + facId + ";").Scan(&faculty_id, &name, &description, &site, &email, &adress, &phone, &university_id)
	if err != nil {
		log.Fatal(err)
	}

	fac := Faculty{
		FacultyId: faculty_id,
		Name: name,
		Description: description,
		Site: site,
		Email: email,
		Adress: adress,
		Phone: phone,
		UniversityId: university_id,
	}

	return fac
}

func getFindUnisNumFromDb(query string) int {
	return getCountFromDb("university_name_descr_vector WHERE name_descr_vector @@ plainto_tsquery('" + query + "')")
}

func getUnisIdsNNamesFromDb(query string) []*University {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var unis []*University
	for rows.Next() {
		var university_id int
		var name string
		err := rows.Scan(&university_id, &name)
		if err != nil {
			log.Fatal(err)
		}

		uni := &University{
			UniversityId: university_id,
			Name: name,
		}

		unis = append(unis, uni)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return unis
}

func findUnisInDb(query string, offset string) []*University {
	dbQuery := "SELECT u.university_id, u.name FROM university u " +
		"JOIN (" +
		"SELECT university_id, ts_rank(name_descr_vector, plainto_tsquery('" + query + "')) " +
		"FROM university_name_descr_vector " +
		"WHERE name_descr_vector @@ plainto_tsquery('" + query + "') " +
		"ORDER BY ts_rank(name_descr_vector, plainto_tsquery('" + query + "')) DESC " +
		"LIMIT 5 OFFSET " + offset +
		") l ON (u.university_id = l.university_id);"

	return getUnisIdsNNamesFromDb(dbQuery)
}

func getUniProfsNumFromDb(uniId string) int {
	from := "(" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		") l"
	return getCountFromDb(from)
}

func getProfsFromDb(query string) []*Profile {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var profs []*Profile
	for rows.Next() {
		var profile_id int
		var name string
		err := rows.Scan(&profile_id, &name)
		if err != nil {
			log.Fatal(err)
		}

		prof := &Profile{
			ProfileId: profile_id,
			Name: name,
		}

		profs = append(profs, prof)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return profs
}

func getUniProfsPageFromDb(uniId string, offset string) []*Profile {
	query := "SELECT p.* FROM profile p " +
		"JOIN (" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		") l ON (p.profile_id = l.profile_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getProfsFromDb(query)
}

func getFacProfsNumFromDb(facId string) int {
	from := "(" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		") l"
	return getCountFromDb(from)
}

func getFacProfsPageFromDb(facId string, offset string) []*Profile {
	query := "SELECT p.* FROM profile p " +
		"JOIN (" +
		"SELECT DISTINCT s.profile_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		") l ON (p.profile_id = l.profile_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getProfsFromDb(query)
}

func getProfFromDb(profId string) Profile {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var profile_id int
	var name string
	err = db.QueryRow("SELECT * FROM profile WHERE profile_id = " + profId + ";").Scan(&profile_id, &name)
	if err != nil {
		log.Fatal(err)
	}

	prof := Profile{
		ProfileId: profile_id,
		Name: name,
	}

	return prof
}

func getSpecsFromDb(query string) []*Speciality {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var specs []*Speciality
	for rows.Next() {
		var speciality_id, profile_id int
		var name string
		var bachelor bool
		err := rows.Scan(&speciality_id, &name, &bachelor, &profile_id)
		if err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	return specs
}

func getUniSpecsNumFromDb(uniId string, profId string) int {
	from := "(" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND s.profile_id = " + profId +
		") l"
	return getCountFromDb(from)
}

func getUniSpecsPageFromDb(uniId string, profId string, offset string) []*Speciality {
	query := "SELECT s.* FROM speciality s " +
		"JOIN (" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND s.profile_id = " + profId +
		") l ON (s.speciality_id = l.speciality_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getSpecsFromDb(query)
}

func getFacSpecsNumFromDb(facId string, profId string) int {
	from := "(" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND s.profile_id = " + profId +
		") l"
	return getCountFromDb(from)
}

func getFacSpecsPageFromDb(facId string, profId string, offset string) []*Speciality {
	query := "SELECT s.* FROM speciality s " +
		"JOIN (" +
		"SELECT DISTINCT s.speciality_id FROM speciality s " +
		"JOIN program pr ON (s.speciality_id = pr.speciality_id) " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND s.profile_id = " + profId +
		") l ON (s.speciality_id = l.speciality_id) " +
		"LIMIT 5 OFFSET " + offset + ";"
	return getSpecsFromDb(query)
}

func getSpecFromDb(specId string) Speciality {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	var speciality_id, profile_id int
	var name string
	var bachelor bool
	err = db.QueryRow("SELECT * FROM speciality WHERE speciality_id = " + specId + ";").Scan(&speciality_id, &name, &bachelor, &profile_id)
	if err != nil {
		log.Fatal(err)
	}

	spec := Speciality{
		SpecialityId: speciality_id,
		Name: name,
		Bachelor: bachelor,
		ProfileId: profile_id,
	}

	return spec
}

func getProgsInfoFromDb(query string) []*Program {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Couldn't connect to db")
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var progs []*Program
	for rows.Next() {
		var program_id uuid.UUID
		var speciality_id int
		var name string
		err := rows.Scan(&program_id, &name, &speciality_id)
		if err != nil {
			log.Fatal(err)
		}

		prog := &Program{
			ProgramId: program_id,
			Name: name,
			SpecialityId: speciality_id,
		}

		progs = append(progs, prog)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return progs
}

func getUniProgsNumFromDb(uniId string) int {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId
	return getCountFromDb(from)
}

func getUniProgsPageFromDb(uniId string, offset string) []*Program {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id FROM program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(query)
}

func getFacProgsNumFromDb(facId string) int {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId
	return getCountFromDb(from)
}

func getFacProgsPageFromDb(facId string, offset string) []*Program {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id FROM program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(query)
}

func getUniSpecProgsNumFromDb(uniId string, specId string) int {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND pr.speciality_id = " + specId
	return getCountFromDb(from)
}

func getUniSpecProgsPageFromDb(uniId string, specId string, offset string) []*Program {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id FROM program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"JOIN university u ON (f.university_id = u.university_id) " +
		"WHERE u.university_id = " + uniId + " AND pr.speciality_id = " + specId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(query)
}

func getFacSpecProgsNumFromDb(facId string, specId string) int {
	from := "program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND pr.speciality_id = " + specId
	return getCountFromDb(from)
}

func getFacSpecProgsPageFromDb(facId string, specId string, offset string) []*Program {
	query := "SELECT pr.program_id, pr.name, pr.speciality_id FROM program pr " +
		"JOIN faculty f ON (pr.faculty_id = f.faculty_id) " +
		"WHERE f.faculty_id = " + facId + " AND pr.speciality_id = " + specId +
		" LIMIT 5 OFFSET " + offset + ";"
	return getProgsInfoFromDb(query)
}