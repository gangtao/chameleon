package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	handlers "chameleon/handlers"

	_ "chameleon/docs"
)

// @title chameleon
// @version 1.0
// @description chameleon is a data stream generator.

// @contact.name Gang Tao
// @contact.url
// @contact.email gang.tao@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http

func main() {
	// Gin instance
	r := gin.New()

	// Routes
	r.GET("/", handlers.HealthCheck)

	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Start server
	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
