package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aufi/must-gather-rest-wrapper/pkg/backend"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	db = backend.ConnectDB("file::memory:?cache=shared")
	router = setupRouter()
}

func TestTriggerGathering(t *testing.T) {
	w := httptest.NewRecorder()
	os.Setenv("TIMEOUT", "123m") // Setting default value for Timeout
	req, _ := http.NewRequest("POST", "/must-gather", strings.NewReader("{\"image\": \"foo.io/bar/image\"}"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Error triggering must-gather request %d", w.Code)
	}

	var gathering *backend.Gathering
	if db.Last(&gathering); gathering.ID == 0 || gathering.Image != "foo.io/bar/image" || gathering.Timeout != "123m" {
		t.Errorf("Cannot find correct Gathering in Database after Create request: %v", gathering)
	}

	// po 10ms zkouset jestli je inprogress? nebo ze se pustila goroutine?
}

func TestGetGathering(t *testing.T) {
	var gathering backend.Gathering
	gathering.ID = 123
	gathering.Status = "inprogress"
	gathering.Image = "foo.io/bar/some_image"
	db.Create(&gathering)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/must-gather/123", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Error getting must-gather execution 123 with: %d", w.Code)
	}
}

func TestListGatherings(t *testing.T) {
	var gathering backend.Gathering
	gathering.ID = 456
	gathering.Status = "completed"
	gathering.Image = "foo.io/bar/some_image"
	db.Create(&gathering)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/must-gather", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Error listing must-gather executions with: %d", w.Code)
	}
}
