package http

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	router := gin.New()

	router.Use(Cors())

	router.POST("/", Run)
	router.POST("/build", Build)
	router.POST("/exec", Exec)
	router.POST("/delete", Delete)

	return router
}
