package backend

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"gorm.io/gorm"
)

func MustGatherExec(gathering *Gathering, db *gorm.DB, archiveFilename string) {
	log.Printf("Starting Must-gather execution #%d", gathering.ID)
	gathering.Status = "inprogress"
	db.Save(&gathering)

	// Prepare destination directory
	dest_directory := gatheringDir(gathering.ID)
	err := os.Mkdir(dest_directory, 0750)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare must-gather command for execution
	cmd := exec.Command("oc")

	// Minimal set of args
	args := []string{"oc", "login", "--token", gathering.AuthToken, "&&", "oc", "adm", "must-gather", "--dest-dir", dest_directory}

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
	//args = append(args, "2>&1") // Reading both stdout&stderr outputs to a single buffer
	log.Printf("Must-gather execution #%d command args: %v", gathering.ID, args)
	cmd.Args = args

	// Execute the must-gather and capture output
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	lineStream := bufio.NewScanner(stdout)
	for lineStream.Scan() {
		gathering.ExecOutput = gathering.ExecOutput + "\n" + lineStream.Text()
		db.Save(&gathering)
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error executing oc adm must-gather command: %v", err)
		gathering.Status = "error"
	}

	// Identify archive file
	cmd = exec.Command("find", dest_directory, "-name", archiveFilename)
	// Expecting a single file with name given by forklift/crane must-gather, might be needed to handle multiple files later (pack all files in dir)
	gatheredArchivePath, err := cmd.Output()
	if err != nil || fmt.Sprintf("%s", gatheredArchivePath) == "" {
		log.Printf("Error finding must-gather result archive: %v", err)
		gathering.Status = "error"
	} else {
		gathering.ArchivePath = strings.TrimSuffix(fmt.Sprintf("%s", gatheredArchivePath), "\n")
		fileInfo, err := os.Stat(gathering.ArchivePath)
		if err != nil {
			log.Printf("Error checking must-gather result archive: %v", err)
			gathering.Status = "error"
		} else {
			gathering.ArchiveSize = uint(fileInfo.Size())
		}

	}

	// Store console output and archive
	if gathering.Status == "inprogress" {
		// Determine user-friendly archive name
		if gathering.CustomName != "" {
			fnameParts := strings.SplitN(path.Base(gathering.ArchivePath), ".", 2) // Add CustomName before first dot in the archive filename
			gathering.ArchiveName = fmt.Sprintf("%s-%s.%s", fnameParts[0], gathering.CustomName, fnameParts[1])
		} else {
			gathering.ArchiveName = path.Base(gathering.ArchivePath)
		}

		gathering.Status = "completed"
	}
	log.Printf("Must-gather execution #%d finished with status: %s", gathering.ID, gathering.Status)
	db.Save(&gathering)
}

func gatheringDir(gatheringID uint) string {
	return fmt.Sprintf("/tmp/must-gather-result-%d", gatheringID)
}
