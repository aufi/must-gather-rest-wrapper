package backend

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gorm.io/gorm"
)

func MustGatherExec(gathering *Gathering, db *gorm.DB) {
	log.Printf("Starting Must-gather execution #%d", gathering.ID)
	gathering.Status = "inprogress"
	db.Save(&gathering)

	// Prepare destination directory
	dest_directory := fmt.Sprintf("/tmp/must-gather-result-%d", gathering.ID)
	err := os.Mkdir(dest_directory, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare must-gather command for execution
	cmd := exec.Command("oc", "adm", "must-gather", "--dest-dir", dest_directory, "--image", "quay.io/konveyor/forklift-must-gather") // + additional args (sanitized to not concat commands)

	// Execute the must-gather
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing oc adm must-gather command: %v", err)
		gathering.Status = "error"
	}

	// Identify archive file
	cmd = exec.Command("find", dest_directory, "-name", "must-gather.tar.gz") // TODO: expected single file with name given by forklift/crane must-gather
	gatheredArchivePath, err := cmd.Output()
	if err != nil || fmt.Sprintf("%s", gatheredArchivePath) == "" {
		log.Printf("Error finding must-gather result archive: %v", err)
		gathering.Status = "error"
	} else {
		gathering.ArchivePath = fmt.Sprintf("%s", gatheredArchivePath)
	}

	// Store console output and archive
	gathering.ExecOutput = fmt.Sprintf("%s", output)
	if gathering.Status == "inprogress" {
		gathering.Status = "completed"
	}
	log.Printf("Must-gather execution #%d finished with status: %s", gathering.ID, gathering.Status)
	db.Save(&gathering)
}
