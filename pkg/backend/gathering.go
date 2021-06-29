package backend

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
	cmd := exec.Command("oc")

	// Minimal set of args
	args := []string{"oc", "adm", "must-gather", "--dest-dir", dest_directory}

	// Expand args for given options (a shared function would need use reflection or marshaling which didn't look to be reasonable to me)
	// ? args sanitized to not concat commands like image="quay.io/foo/bar; rm -rf something"
	if gathering.Image != "" {
		args = append(args, "--image", gathering.Image)
	}
	if gathering.ImageStream != "" {
		args = append(args, "--image-stream", gathering.ImageStream)
	}
	if gathering.NodeName != "" {
		args = append(args, "--node-name", gathering.NodeName)
	}
	if gathering.SourceDir != "" {
		args = append(args, "--source-dir", gathering.SourceDir)
	}
	if gathering.Timeout != "" {
		args = append(args, "--timeout", gathering.Timeout)
	}
	if gathering.Server != "" {
		args = append(args, "--server", gathering.Server)
	}
	if gathering.Command != "" {
		args = append(args, "--", gathering.Command)
	}
	log.Printf("Must-gather execution #%d command args: %v", gathering.ID, args)
	cmd.Args = args

	// Execute the must-gather
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing oc adm must-gather command: %v", err)
		gathering.Status = "error"
	}

	// Identify archive file
	cmd = exec.Command("find", dest_directory, "-name", "must-gather.tar.gz")
	// Expecting a single file with name given by forklift/crane must-gather, might be needed to handle multiple files later (pack all files in dir)
	gatheredArchivePath, err := cmd.Output()
	if err != nil || fmt.Sprintf("%s", gatheredArchivePath) == "" {
		log.Printf("Error finding must-gather result archive: %v", err)
		gathering.Status = "error"
	} else {
		gathering.ArchivePath = strings.TrimSuffix(fmt.Sprintf("%s", gatheredArchivePath), "\n")
	}

	// Store console output and archive
	gathering.ExecOutput = fmt.Sprintf("%s", output)
	if gathering.Status == "inprogress" {
		gathering.Status = "completed"
	}
	log.Printf("Must-gather execution #%d finished with status: %s", gathering.ID, gathering.Status)
	db.Save(&gathering)
}
