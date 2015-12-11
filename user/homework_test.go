package user

import (
	"testing"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

func TestRejectMissingtracks(t *testing.T) {
	problem1 := &api.Problem{
		TrackID:  "go",
		Slug:     "clock",
		Language: "Go",
		Name:     "Clock",
	}
	problem2 := &api.Problem{
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

	hw := NewHomework([]*api.Problem{problem1, problem2}, &config.Config{Dir: "/tmp"})
	err := hw.RejectMissingTracks(dirMap)

	if err != nil {
		t.Error(err)
	}

	if len(hw.Items) == 2 {
		t.Error("Should have rejected the Ruby problem but did not reject any problems!")
	}

	if len(hw.Items) == 1 && hw.Items[0].TrackID == "ruby" {
		t.Error("Should have rejected the Ruby problem and rejected the Go problem instead!")
	}

	if err := hw.RejectMissingTracks(emptyDirMap); err == nil {
		t.Error("Should have returned error because user hasn't started any tracks but didn't!")
	}
}
