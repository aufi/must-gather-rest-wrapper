package backend

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Gathering struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created-at"`
	UpdatedAt   time.Time `json:"updated-at"`
	AuthToken   string    `json:"-"` // Maybe use hash function to not store the real token
	CustomName  string    `gorm:"index" form:"custom-name" json:"custom-name"`
	Status      string    `json:"status"` // Expected values: new, inprogress, completed, error
	Image       string    `form:"image" json:"image"`
	ImageStream string    `form:"image-stream" json:"image-stream"`
	NodeName    string    `form:"node-name" json:"node-name"`
	Command     string    `form:"command" json:"command"`
	SourceDir   string    `form:"source-dir" json:"source-dir"`
	Timeout     string    `form:"timeout" json:"timeout"`
	Server      string    `form:"server" json:"server"`
	ArchivePath string    `json:"-"` // Not exposed via JSON API
	ArchiveSize uint      `json:"archive-size"`
	ArchiveName string    `json:"archive-name"`
	ExecOutput  string    `json:"exec-output"` // Fields without form:"<name>" cannot be set via API by bind
}

func ConnectDB(databasePath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databasePath))
	if err != nil {
		panic("Failed to connect database")
	}

	err = db.AutoMigrate(&Gathering{})
	if err != nil {
		panic("Failed to migrate database")
	}
	return db
}

func PeriodicalCleanup(maxAgeOption string, db *gorm.DB, runOnceNow bool) {
	if maxAgeOption == "" || maxAgeOption == "-1" {
		log.Println("Periodical must-gather executions cleanup is disabled.")
		return
	}
	age, err := time.ParseDuration(maxAgeOption)
	if err != nil {
		log.Panic(fmt.Printf("Error parsing maxAgeOption for cleanup: %s", err))
	}
	log.Printf("Scheduling must-gather executions cleanup for older than %s", age)
	tickerPeriod := 1 * time.Hour
	if runOnceNow == true {
		tickerPeriod = 1 * time.Millisecond
	}
	for currTime := range time.Tick(tickerPeriod) {
		oldTime := currTime.Add(-1 * age)
		log.Printf("Checking outdated must-gather executions")
		var obsoleteGatherings []*Gathering
		db.Where("updated_at < ?", oldTime).Find(&obsoleteGatherings)
		for _, obsoleteGathering := range obsoleteGatherings {
			log.Printf("Deleting outdated must-gather execution #%d", obsoleteGathering.ID)
			// Remove gathering data directory
			err := os.RemoveAll(gatheringDir(obsoleteGathering.ID))
			if err != nil {
				log.Printf("Error delete directory: %v", err)
			}
			// Remove gathering record from database
			db.Delete(&obsoleteGathering)
		}
		if runOnceNow == true {
			break
		}
	}
}
