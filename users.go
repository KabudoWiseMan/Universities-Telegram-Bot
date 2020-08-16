package main

import "math"

const (
	NoState = 0
	UniState = 1
	FindUniState = 2
	RatingQSState = 3
	FeeState = 4
	CityState = 5
	ProfileState = 6
	SpecialityState = 7
	EgeState = 8
	SubjState = 9
	ProgState = 10
)

type Ege struct {
	SubjId int
	MinPoints uint64
}

type UserInfo struct {
	State int
	Query string
	City int
	Dormatary bool
	MilitaryDep bool
	Fee uint64
	ProfileId int
	SpecialityId int
	EntryTest bool
	Eges []Ege
	LastSubj int
}

func (usr *UserInfo) Clear() {
	usr.City = 0
	usr.Dormatary = false
	usr.MilitaryDep = false
	usr.Fee = math.MaxUint64
	usr.ProfileId = 0
	usr.SpecialityId = 0
	usr.EntryTest = false
	usr.Eges = nil
}

func (usr *UserInfo) DeleteEge() {
	var newEges []Ege
	for _, ege := range usr.Eges {
		if ege.SubjId != usr.LastSubj {
			newEges = append(newEges, ege)
		}
	}

	usr.Eges = newEges
}

func (usr *UserInfo) AddEge(points uint64) bool {
	found := false
	var newEges []Ege
	for _, ege := range usr.Eges {
		if ege.SubjId == usr.LastSubj {
			ege.MinPoints = points
			found = true
		}
		newEges = append(newEges, ege)
	}

	if !found {
		newEges = append(newEges, Ege{usr.LastSubj, points})
	}

	usr.Eges = newEges

	return found
}

type Users struct {
	users map[int64]*UserInfo
}

func InitUsers() *Users {
	return &Users{make(map[int64]*UserInfo)}
}

func (usrs *Users) User(userId int64) *UserInfo {
	if user, ok := usrs.users[userId]; ok {
		return user
	} else {
		newUser := &UserInfo {
			State: NoState,
			Fee: math.MaxUint64,
		}

		usrs.users[userId] = newUser

		return newUser
	}
}

func (usrs *Users) Delete(userId int64) {
	delete(usrs.users, userId)
}
