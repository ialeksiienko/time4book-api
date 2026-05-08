package main

import (
	_ "time4book/docs"
	"time4book/internal"
)

// @title           Time4Book
// @version         1.0
// @description     Time4Book service API.
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:50052
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app := internal.New()

	server := app.App.GetHTTPServer()

	server.Run()
}
