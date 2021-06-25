package backend

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Gathering struct {
	gorm.Model
	Status      string `form:"status" json:"status"` // new, inprogress, completed, error
	Image       string `form:"image" json:"image"`
	ImageStream string `form:"image-stream" json:"image-stream"`
	NodeName    string `form:"node-name" json:"node-name"`
	Command     string `form:"command" json:"command"`
	SourceDir   string `form:"source-dir" json:"source-dir"`
	Timeout     string `form:"timeout" json:"timeout"`
	ArchivePath string
	//ArchiveUrl    string
	ExecOutput string
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
