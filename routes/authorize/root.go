package authorize

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(AuthorizeGroup{})

type AuthorizeGroup struct {
	courier.EmptyOperator
}

func (AuthorizeGroup) Path() string {
	return "/authorize"
}
