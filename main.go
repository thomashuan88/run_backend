package main

import (
	"run-backend/tasks"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	tasks.SetupGameboRoutes(r)
	r.Run(":8080")
}
