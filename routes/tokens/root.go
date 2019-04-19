package tokens

import "github.com/johnnyeven/libtools/courier"

var Router = courier.NewRouter(TokenGroup{})

type TokenGroup struct {
	courier.EmptyOperator
}

func (TokenGroup) Path() string {
	return "/tokens"
}
