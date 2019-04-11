package main

import (
	"github.com/johnnyeven/libtools/servicex"

	"github.com/in2store/service-in2-auth/global"
	"github.com/in2store/service-in2-auth/routes"
)

func main() {
	servicex.Execute()
	global.Config.Server.Serve(routes.RootRouter)
}
