package routes

import (
	"github.com/in2store/service-in2-auth/routes/authorize"
	"github.com/in2store/service-in2-auth/routes/channels"
	"github.com/in2store/service-in2-auth/routes/sessions"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/swagger"
)

var RootRouter = courier.NewRouter(GroupRoot{})
var V0Router = courier.NewRouter(V0Group{})

func init() {
	RootRouter.Register(swagger.SwaggerRouter)
	RootRouter.Register(V0Router)

	V0Router.Register(channels.Router)
	V0Router.Register(authorize.Router)
	V0Router.Register(sessions.Router)
}

type GroupRoot struct {
	courier.EmptyOperator
}

func (root GroupRoot) Path() string {
	return "/in2-auth"
}

type V0Group struct {
	courier.EmptyOperator
}

func (V0Group) Path() string {
	return "/v0"
}
