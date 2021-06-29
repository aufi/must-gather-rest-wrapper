package backend

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Gathering struct {
	ID          uint   `gorm:"primarykey"`
	Status      string `json:"status"` // Expected values: new, inprogress, completed, error
	Image       string `form:"image" json:"image"`
	ImageStream string `form:"image-stream" json:"image-stream"`
	NodeName    string `form:"node-name" json:"node-name"`
	Command     string `form:"command" json:"command"`
	SourceDir   string `form:"source-dir" json:"source-dir"`
	Timeout     string `form:"timeout" json:"timeout"`
	Server      string `form:"server" json:"server"`
	ArchivePath string `json:"-"`           // Not exposed via JSON API
	ExecOutput  string `json:"exec-output"` // Fields without form:"<name>" cannot be set via API by bind
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
