package main

import (
	_ "go.bankyaya.org/app/backend/cmd/swagger/docs"
	"go.bankyaya.org/app/backend/internal/adapter/http/server"
	"go.bankyaya.org/app/backend/internal/pkg/config"
)

type app struct {
	ss *server.Server
}

func newApp(ss *server.Server) *app {
	return &app{
		ss: ss,
	}
}

// main swaggo annotation.
//
//	@title			API Specification
//	@version		1.0
//	@description	Greet service API specification.
//	@termsOfService	https://swagger.io/terms/
//	@contact.name	BillyKore
//	@contact.url	https://www.swagger.io/support
//	@contact.email	billyimmcul2010@gmail.com
//	@license.name	Apache 2.0
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0.html
//	@host			api.bankyaya.co.id
//	@schemes		http https
//	@BasePath		/api/v1
func main() {
	c := config.Load()
	a := initApp(c)

	a.ss.Serve()
}
