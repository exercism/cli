package user

import (
	"testing"

	"github.com/robphoenix/cli/api"
	"github.com/robphoenix/cli/config"
)

func TestRejectMissingtracks(t *testing.T) {
	exercise1 := &api.Exercise{
		TrackID:  "go",
		Slug:     "clock",
		Language: "Go",
		Name:     "Clock",
	}
	exercise2 := &api.Exercise{
		TrackID:  "ruby",
		Slug:     "clock",
		Language: "Ruby",
		Name:     "Clock",
	}
	dirMap := map[string]bool{
		"/tmp/go":      true,
		"/tmp/haskell": true,
	}
	emptyDirMap := make(map[string]bool)

	hw := NewHomework([]*api.Exercise{exercise1, exercise2}, &config.Config{Dir: "/tmp"})
	err := hw.RejectMissingTracks(dirMap)

	if err != nil {
		t.Error(err)
	}

	if len(hw.Items) == 2 {
		t.Error("Should have rejected the Ruby exercise but did not reject any exercises!")
	}

	if len(hw.Items) == 1 && hw.Items[0].TrackID == "ruby" {
		t.Error("Should have rejected the Ruby exercise and rejected the Go exercise instead!")
	}

	if err := hw.RejectMissingTracks(emptyDirMap); err == nil {
		t.Error("Should have returned error because user hasn't started any tracks but didn't!")
	}
}
