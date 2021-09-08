package main

import (
	"log"
	"os"
	"strings"

	"github.com/aufi/must-gather-rest-wrapper/pkg/backend"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var r *gin.Engine
var db *gorm.DB

func main() {
	db = backend.ConnectDB(configEnvOrDefault("DB_PATH", "gatherings.db"))

	// Periodical deletion of old records&archives on background
	go backend.PeriodicalCleanup(configEnvOrDefault("CLEANUP_MAX_AGE", "-1"), db, false)

	// Start HTTP service
	r := setupRouter()
	r.Run() // PORT from ENV variable is handled inside Gin-gonic and defaults to 8080
}

func setupRouter() *gin.Engine {
	// Gin routes setup
	r = gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "must-gather-rest-wrapper - for API see https://github.com/aufi/must-gather-rest-wrapper/tree/main/doc")
	})

	r.POST("/must-gather", triggerGathering)
	r.GET("/must-gather", listGatherings) // good at least during development and testing, real user should know gathering ID
	r.GET("/must-gather/:id", getGathering)
	r.GET("/must-gather/:id/data", getGatheringArchive)

	return r
}

func triggerGathering(c *gin.Context) {
	var gathering backend.Gathering
	if err := c.Bind(&gathering); err == nil {
		gathering.Status = "new"
		if gathering.Image == "" {
			gathering.Image = configEnvOrDefault("DEFAULT_IMAGE", "quay.io/konveyor/forklift-must-gather") // default image configurable via OS ENV vars
		}
		if gathering.Timeout == "" {
			gathering.Timeout = configEnvOrDefault("TIMEOUT", "20m") // default timeout for must-gather execution
		}
		gathering.AuthToken = strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")
		// TODO: Check the token or just pass to the commandline? but always satinize to not explode token to multiple commands (steal previous executions data or tokens)
		db.Create(&gathering)
		c.JSON(201, gathering)
		go backend.MustGatherExec(&gathering, db, configEnvOrDefault("ARCHIVE_FILENAME", "must-gather.tar.gz"))
	} else {
		log.Printf("Error creating gathering: %v", err)
		c.JSON(400, "create gathering error")
	}
}

// TODO filter results by authtoken
func listGatherings(c *gin.Context) {
	var Gatherings []backend.Gathering
	db.Find(&Gatherings)
	c.JSON(200, Gatherings)
}

func getGathering(c *gin.Context) {
	var gathering backend.Gathering
	db.First(&gathering, "id = ?", c.Param("id")) // ID (uint) lookup - safe way to handle possible string input to not interpret it as a query
	if gathering.ID != 0 {
		c.JSON(200, gathering)
	} else {
		db.Last(&gathering, "custom_name = ?", c.Param("id")) // Fallback to CustomName (string) lookup - returned the newest/last matching record
		if gathering.ID != 0 {
			c.JSON(200, gathering)
		} else {
			// Return empty gathering with 404 code if not found
			c.JSON(404, gathering)
		}
	}
}

func getGatheringArchive(c *gin.Context) {
	var gathering backend.Gathering
	db.First(&gathering, "id = ?", c.Param("id"))
	if gathering.ID != 0 && gathering.Status == "completed" {
		c.FileAttachment(gathering.ArchivePath, gathering.ArchiveName)
	} else {
		c.String(404, "")
	}
}

func configEnvOrDefault(name, defaultValue string) string {
	value, present := os.LookupEnv(name)
	if present {
		log.Printf("Config option %s set from environment variable to: %s", name, value)
		return value
	} else {
		log.Printf("Environment variable %s is undefined, using default: %s", name, defaultValue)
		return defaultValue
	}
}
