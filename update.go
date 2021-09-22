package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func uniIsWrong(uni *University) bool {
	return uni == nil ||
		uni.Name == "" && uni.Description == "" ||
		uni.Site == "" && uni.Email == "" && uni.Address == "" && uni.Phone == ""
}

func facIsWrong(fac *Faculty) bool {
	return fac == nil ||
		fac.Name == "" && fac.Description == "" ||
		fac.Site == "" && fac.Email == "" && fac.Address == "" && fac.Phone == ""
}

func progIsWrong(prog *Program) bool {
	return prog == nil ||
		prog.Name == "" ||
		prog.SpecialityId == -1 ||
		prog.FreePassPoints == -1 ||
		prog.FreePlaces == -1 ||
		prog.PaidPlaces == -1 ||
		prog.Fee == -1 ||
		prog.PaidPassPoints == -1 ||
		prog.StudyForm == "" && prog.StudyLanguage == "" && prog.StudyBase == "" && prog.StudyYears == "" && prog.Description == ""
}

func minPointsIsWrong(minPoints *MinEgePoints) bool {
	return minPoints == nil ||
		minPoints.SubjectId == 0
}

func entrTestIsWrong(entrTest *EntranceTest) bool {
	return entrTest == nil ||
		entrTest.TestName == "" ||
		entrTest.MinPoints == 0
}

func ratingQsIsWrong(ratingQS *RatingQS) bool {
	return ratingQS == nil ||
		ratingQS.UniversityId == -1 ||
		ratingQS.HighMark == 0 ||
		ratingQS.LowMark == 0
}

func profIsWrong(prof *Profile) bool {
	return prof == nil ||
		prof.ProfileId == 0 ||
		prof.Name == ""
}

func specIsWrong(spec *Speciality) bool {
	return spec == nil ||
		spec.SpecialityId == 0 ||
		spec.Name == ""
}

func parseAndUpdateUnis(db *sql.DB) error {
	cities, err := getAllCitiesFromDb(db)
	if err != nil {
		log.Println("couldn't get cities from db for update universities, error:", err)
		return err
	}

	log.Println("Parsing universities started")
	startParse := time.Now()
	unis := parseUniversities(cities)
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	if len(unis) == 0 || uniIsWrong(unis[0]) {
		log.Println("Parsing universities failed")
		return errors.New("error")
	}
	log.Println("Parsing universities finished")

	log.Println("Updating universities started")
	startUpdate := time.Now()
	if err := updateUnisInDb(db, unis); err != nil {
		log.Println("Updating universities failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating universities finished")
	return nil
}

func updateUnis() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update universities", err)
		return
	}
	log.Println("Successfully connected to data base for update universities")
	defer closeDb(db)

	parseAndUpdateUnis(db)
}

func parseAndUpdateFacs(db *sql.DB) error {
	unis, err := getUnisIdsNamesFromDb(db, false)
	if err != nil {
		log.Println("couldn't get universities from db for update faculties, error:", err)
		return err
	}

	log.Println("Parsing faculties started")
	startParse := time.Now()
	facs := parseFaculties(unis)
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	if len(facs) == 0 || facIsWrong(facs[0]) {
		log.Println("Parsing faculties failed")
		return errors.New("error")
	}
	log.Println("Parsing faculties finished")

	log.Println("Updating faculties started")
	startUpdate := time.Now()
	if err := updateFacsInDb(db, facs); err != nil {
		log.Println("Updating faculties failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating faculties finished")
	return nil
}

func updateFacs() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update faculties", err)
		return
	}
	log.Println("Successfully connected to data base for update faculties")
	defer closeDb(db)

	parseAndUpdateFacs(db)
}

func parseAndUpdateProgsNInfo(db *sql.DB) error {
	facs, err := getFacsIdsFromDb(db)
	if err != nil {
		log.Println("couldn't get faculties from db for update programs, error:", err)
		return err
	}

	var specs []*Speciality
	specs, err = getSpecsIdsFromDb(db)
	if err != nil {
		log.Println("couldn't get specialities from db for update programs, error:", err)
		return err
	}

	var subjs map[string]int
	subjs, err = getRevSubjsMapFromDb(db)
	if err != nil {
		log.Println("couldn't get subjects from db for update programs, error:", err)
		return err
	}

	if err = createTempTablesForProgsNInfo(db); err != nil {
		log.Println("couldn't create temp tables for programs, error:", err)
		return err
	}

	pace := 100

	log.Println("Parsing and inserting programs started")
	startParse := time.Now()
	for i := 0; i < len(facs); i += pace {
		to := i + pace
		if to > len(facs) {
			to = len(facs)
		}

		facsSlice := facs[i : to]

		log.Println("Parsing and inserting programs of faculties from", i, "to", to, "started")
		for j := 0; j <= 3; j++ {
			progs, minPoints, entrTests := parsePrograms(facsSlice, specs, subjs)
			if len(progs) == 0 || len(minPoints) == 0 || len(entrTests) == 0 ||
				progIsWrong(progs[0]) || minPointsIsWrong(minPoints[0]) || entrTestIsWrong(entrTests[0]) {
				err = errors.New("parsing error")
			} else {
				err = insertProgsNInfoToTempTables(db, progs, minPoints, entrTests)
				if err == nil {
					break
				}
			}
		}
		if err != nil {
			log.Println("Parsing and inserting programs of faculties from", i, "to", to, "failed with error:", err)
			return err
		}
		log.Println("Parsing and inserting programs of faculties from", i, "to", to, "finished")
	}
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	log.Println("Parsing and inserting programs finished")

	log.Println("Updating programs started")
	startUpdate := time.Now()
	if err := updateProgsNInfoInDb(db); err != nil {
		log.Println("Updating programs failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating programs finished")

	return nil
}

func updateProgsNInfo() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update programs", err)
		return
	}
	log.Println("Successfully connected to data base for update programs")
	defer closeDb(db)

	parseAndUpdateProgsNInfo(db)
}

func parseAndUpdateCities(db *sql.DB) error {
	log.Println("Parsing cities started")
	startParse := time.Now()
	cities := parseCities()
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	if len(cities) == 0 || getElemFromMap(cities) == "" {
		log.Println("Parsing cities failed")
		return errors.New("error")
	}
	log.Println("Parsing cities finished")

	log.Println("Updating cities started")
	startUpdate := time.Now()
	if err := updateCitiesInDb(db, cities); err != nil {
		log.Println("Updating cities failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating cities finished")
	return nil
}

func updateCities() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update cities", err)
		return
	}
	log.Println("Successfully connected to data base for update cities")
	defer closeDb(db)

	parseAndUpdateCities(db)
}

func parseAndUpdateSubjs(db *sql.DB) error {
	log.Println("Parsing subjects started")
	startParse := time.Now()
	subjs := parseSubjs()
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	if len(subjs) == 0 {
		log.Println("Parsing subjects failed")
		return errors.New("error")
	}
	if _, ok := subjs[""]; ok {
		log.Println("Parsing subjects failed")
		return errors.New("error")
	}
	log.Println("Parsing subjects finished")

	log.Println("Updating subjects started")
	startUpdate := time.Now()
	if err := updateSubjsInDb(db, subjs); err != nil {
		log.Println("Updating subjects failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating subjects finished")
	return nil
}

func updateSubjs() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update subjects", err)
		return
	}
	log.Println("Successfully connected to data base for update subjects")
	defer closeDb(db)

	parseAndUpdateSubjs(db)
}

func parseAndUpdateRatingQS(db *sql.DB) error {
	unis, err := getUnisWithSitesFromDb(db)
	if err != nil {
		log.Println("couldn't get universities from db for update rating QS, error:", err)
		return err
	}

	log.Println("Parsing rating QS started")
	startParse := time.Now()
	ratingQS := parseRatingQS(unis)
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)
	if len(ratingQS) == 0 || ratingQsIsWrong(ratingQS[0]) {
		log.Println("Parsing rating QS failed")
		return errors.New("error")
	}
	log.Println("Parsing rating QS finished")

	for _, rate := range ratingQS {
		fmt.Println(*rate)
	}

	log.Println("Updating rating QS started")
	startUpdate := time.Now()
	if err := updateRatingQSInDb(db, ratingQS); err != nil {
		log.Println("Updating rating QS failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating rating QS finished")
	return nil
}

func updateRatingQS() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update rating QS", err)
		return
	}
	log.Println("Successfully connected to data base for update rating QS")
	defer closeDb(db)

	parseAndUpdateRatingQS(db)
}

func parseAndUpdateProfsNSpecs(db *sql.DB) error {
	log.Println("Parsing profiles and specialities started")
	startParse := time.Now()
	profsBach, specs := parseProfsNSpecs(BachelorSpecialitiesSite)
	profsSpec, specsSpec := parseProfsNSpecs(SpecialistSpecialitiesSite)
	stopParse := time.Since(startParse)
	fmt.Println("Parsing time: ", stopParse)

	profsMap := make(map[Profile]bool)
	for _, p := range profsBach {
		profsMap[*p] = true
	}
	for _, p := range profsSpec {
		profsMap[*p] = true
	}
	profsBach = nil
	profsSpec = nil

	var profs []*Profile
	for profM, _ := range profsMap {
		prof := &Profile{
			ProfileId: profM.ProfileId,
			Name: profM.Name,
		}
		profs = append(profs, prof)
	}
	profsMap = nil

	specs = append(specs, specsSpec...)
	specsSpec = nil

	if len(profs) == 0 || len(specs) == 0 || profIsWrong(profs[0]) || specIsWrong(specs[0]) {
		log.Println("Parsing profiles and specialities failed")
		return errors.New("error")
	}
	log.Println("Parsing profiles and specialities finished")

	log.Println("Updating profiles and specialities started")
	startUpdate := time.Now()
	if err := updateProfsNSpecsInDb(db, profs, specs); err != nil {
		log.Println("Updating profiles and specialities failed with error:", err)
		return err
	}
	stopUpdate := time.Since(startUpdate)
	fmt.Println("Update time: ", stopUpdate)
	log.Println("Updating profiles and specialities finished")
	return nil
}

func updateProfsNSpecs() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update profiles and specialities", err)
		return
	}
	log.Println("Successfully connected to data base for update profiles and specialities")
	defer closeDb(db)

	parseAndUpdateProfsNSpecs(db)
}

func updateDb() {
	db, err := connectToDb()
	if err != nil {
		log.Println("couldn't connected to data base for update", err)
		return
	}
	log.Println("Successfully connected to data base for update")
	defer closeDb(db)

	log.Println("Update started")

	fmt.Println("CITIES:")
	if err = parseAndUpdateCities(db); err != nil {
		return
	}

	fmt.Println("UNIS:")
	if err = parseAndUpdateUnis(db); err != nil {
		return
	}

	fmt.Println("FACS:")
	if err = parseAndUpdateFacs(db); err != nil {
		return
	}

	fmt.Println("PROFS & SPECS:")
	if err = parseAndUpdateProfsNSpecs(db); err != nil {
		return
	}

	fmt.Println("SUBJS:")
	if err = parseAndUpdateSubjs(db); err != nil {
		return
	}

	fmt.Println("PROGS & Info:")
	if err = parseAndUpdateProgsNInfo(db); err != nil {
		return
	}

	fmt.Println("RATING:")
	if err = parseAndUpdateRatingQS(db); err != nil {
		return
	}

	log.Println("Update finished")
}