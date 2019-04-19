package sessions

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(SessionGroup{})

type SessionGroup struct {
	courier.EmptyOperator
}

func (SessionGroup) Path() string {
	return "/sessions"
}
