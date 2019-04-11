package routes

import (
	"github.com/in2store/service-in2-auth/routes/authorize"
	"github.com/in2store/service-in2-auth/routes/channel"
	"github.com/johnnyeven/libtools/courier"
	"github.com/johnnyeven/libtools/courier/swagger"
)

var RootRouter = courier.NewRouter(GroupRoot{})
var V0Router = courier.NewRouter(V0Group{})

func init() {
	RootRouter.Register(swagger.SwaggerRouter)
	RootRouter.Register(V0Router)

	V0Router.Register(channel.Router)
	V0Router.Register(authorize.Router)
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
