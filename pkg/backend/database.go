package backend

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Gathering struct {
	gorm.Model
	Status      string `json:"status"` // Expected values: new, inprogress, completed, error
	Image       string `form:"image" json:"image"`
	ImageStream string `form:"image-stream" json:"image-stream"`
	NodeName    string `form:"node-name" json:"node-name"`
	Command     string `form:"command" json:"command"`
	SourceDir   string `form:"source-dir" json:"source-dir"`
	Timeout     string `form:"timeout" json:"timeout"`
	Server      string `form:"server" json:"server"`
	ArchivePath string `json:"archive-path"` // do not expose via json API?
	ExecOutput  string `json:"exec-output"`  // Fields without form:"<name>" cannot be set via API by bind
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
