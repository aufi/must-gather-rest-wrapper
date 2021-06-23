package main

import (
	"log"

	"github.com/aufi/must-gather-rest-wrapper/pkg/backend"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var r *gin.Engine
var db *gorm.DB

func main() {
	db = backend.ConnectDB("gatherings_dev.db") //("file::memory:?cache=shared") // Ephemeral database backend until persistence would be needed
	r = gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "must-gather-rest-wrapper - for API see https://github.com/aufi/must-gather-rest-wrapper/tree/main/doc")
	})

	r.POST("/must-gather", triggerGathering)
	r.GET("/must-gather", listGatherings) // good at least during development and testing, real user should know gathering ID
	r.GET("/must-gather/:id", getGathering)

	r.Run(":8080")
}

func triggerGathering(c *gin.Context) {
	var gathering backend.Gathering
	if err := c.Bind(&gathering); err == nil {
		gathering.Status = "new"
		db.Create(&gathering)
		c.JSON(201, gathering)
		go backend.MustGatherExec(&gathering, db)
	} else {
		log.Printf("Error creating gathering: %v", err)
		c.JSON(400, "create gathering error")
	}
}

func listGatherings(c *gin.Context) {
	var Gatherings []backend.Gathering
	db.Find(&Gatherings)
	c.JSON(200, Gatherings)
}

func getGathering(c *gin.Context) {
	var gathering backend.Gathering
	db.First(&gathering, "id = ?", c.Param("id")) // safe way to handle possible string input to not interpret it as a query
	if gathering.ID != 0 {
		c.JSON(200, gathering)
	} else {
		c.JSON(404, "not found")
	}
}
