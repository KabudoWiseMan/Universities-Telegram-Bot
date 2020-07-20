package main

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