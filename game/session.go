package game

type session struct {
	users map[int]int // userID:connID
}

func (s *session) handleHandshake() {

}

func (s *session) handleHeartbeat() {

}

func (s *session) getUserConnID() {

}
