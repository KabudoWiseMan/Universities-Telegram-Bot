package main

const (
	NoState = 0
	UiState = 1
	FindUniState = 2
	RatingQSState = 3
)

type UserInfo struct {
	State int
	Page int
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
		newUser := &UserInfo{
			State: NoState,
			Page: 1,
		}

		usrs.users[userId] = newUser

		return newUser
	}
}

func (usrs *Users) Delete(userId int64) {
	delete(usrs.users, userId)
}
