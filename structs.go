package main

import uuid "github.com/satori/go.uuid"

type Profile struct {
	ProfileId int
	Name string
}

type Speciality struct {
	SpecialityId int
	Name string
	Bachelor bool
	ProfileId int
}

type University struct {
	UniversityId int
	Name string
	Description string
	Site string
	Email string
	Adress string
	Phone string
	MilitaryDep bool
	Dormitary bool
}

type Faculty struct {
	FacultyId int
	Name string
	Description string
	Site string
	Email string
	Adress string
	Phone string
	UniversityId int
}

type Program struct {
	ProgramId uuid.UUID
	ProgramNum int
	Name string
	Description string
	FreePlaces int
	PaidPlaces int
	Fee float64
	FreePassPoints int
	PaidPassPoints int
	StudyForm string
	StudyLanguage string
	StudyBase string
	StudyYears string
	FacultyId int
	SpecialityId int
}

type MinEgePoints struct {
	ProgramId uuid.UUID
	SubjectId int
	MinPoints int
}

type EntranceTest struct {
	ProgramId uuid.UUID
	TestName string
	MinPoints int
}

