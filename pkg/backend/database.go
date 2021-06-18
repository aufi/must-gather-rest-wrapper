package backend

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Gathering struct {
	gorm.Model
	Status        string `json:"status"` // new, inprogress, completed, error
	Image         string `json:"image" binding:"required"`
	CustomCommand string
	ArchiveLink   string
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
