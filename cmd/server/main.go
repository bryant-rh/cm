package main

import (
	_ "github.com/bryant-rh/cm/cmd/server/docs"
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/cmd/server/routers"

	"github.com/kunlun-qilian/confserver"
	"github.com/spf13/cobra"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	confserver.Execute(func(cmd *cobra.Command, args []string) {
		s := global.Config.Server
		routers.NewRooter(s.Engine())
		s.Serve()
	})
}
