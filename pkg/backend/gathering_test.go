package backend

import (
	"os"
	"testing"
)

func TestMustGatherExec(t *testing.T) {
	// Prepare gathering in a fresh db
	db := ConnectDB("file::memory:?cache=shared")
	var gathering Gathering
	gathering.Status = "new"
	gathering.Image = "foo.bar/image"
	gathering.Timeout = "11m"
	db.Create(&gathering)

	// Exec the gathering
	MustGatherExec(&gathering, db, "must-gather.tar.gz")

	// Check gathering after execution
	db.Last(&gathering, gathering.ID)
	if gathering.Status != "completed" {
		t.Error("Must-gather exection status should be completed")
	}

	if gathering.ExecOutput == "" {
		t.Error("Must-gather exec output should not be empty")
	}

	if _, err := os.Stat(gathering.ArchivePath); err != nil {
		t.Errorf("Archive file is missing, looked for: %s", gathering.ArchivePath)
	}
}
