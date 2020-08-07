package main

const (
	NoState = 0
	RatingQSState = 1
)

type FSM struct {
	users map[int64]*UserInfo
}

type UserInfo struct {
	state int
	page int
}

func (fsm *FSM) GetUser(userId int64) (*UserInfo, bool) {
	if user, ok := fsm.users[userId]; ok {
		return user, true
	}

	return nil, false
}

func (fsm *FSM) State(userId int64) int {
	if user, ok := fsm.users[userId]; ok {
		return user.state
	}

	return 0
}

func (fsm *FSM) Page(userId int64) int {
	if user, ok := fsm.users[userId]; ok {
		return user.page
	}

	return 0
}

func (fsm *FSM) SetState(userId int64, newState int) {
	if user, ok := fsm.users[userId]; ok {
		user.state = newState
	} else {
		newUser := &UserInfo{
			state: newState,
			page: 0,
		}

		fsm.users[userId] = newUser
	}
}

func (fsm *FSM) SetPage(userId int64, newPage int) {
	if user, ok := fsm.users[userId]; ok {
		user.page = newPage
	} else {
		newUser := &UserInfo{
			state: 0,
			page: newPage,
		}

		fsm.users[userId] = newUser
	}
}

func (fsm *FSM) Delete(userId int64) {
	delete(fsm.users, userId)
}
