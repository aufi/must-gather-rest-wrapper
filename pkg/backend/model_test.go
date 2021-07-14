package backend

import (
	"testing"
)

func TestPeriodicalCleanup(t *testing.T) {
	// Prepare gatherings in a fresh db
	db := ConnectDB("file::memory:?cache=shared")
	var freshGathering, oldGathering Gathering
	oldGathering.ID = 123
	oldGathering.Status = "completed"
	freshGathering.ID = 456
	freshGathering.Status = "completed"
	db.Create(&oldGathering)
	db.Create(&freshGathering)
	// Adjust updatedAt field
	db.Model(&oldGathering).UpdateColumn("updated_at", "2021-01-01T00:00:00.0+02:00")

	// Run the cleanup
	PeriodicalCleanup("12h", db, true)

	// Check db after cleanup
	var gathering Gathering
	db.Last(&gathering, oldGathering.ID) // oldGathering should be deleted
	if gathering.ID != 0 {
		t.Errorf("The oldGathering is still present in db with ID: %d", gathering.ID)
	}
	db.Last(&gathering, freshGathering.ID) // freshGathering should be kept
	if gathering.ID == 0 {
		t.Errorf("The freshGathering is not present in db, looked for ID: %d", freshGathering.ID)
	}
}
