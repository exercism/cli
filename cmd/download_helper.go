package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	netURL "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/workspace"
)

type downloadParams struct {
	cfg   *config.Config
	uuid  string
	slug  string
	track string
	team  string
}

type downloadPayloadContext struct {
	cfg     *config.Config
	payload *downloadPayload
}

func newDownloadPayload(params downloadParams) (*downloadPayloadContext, error) {
	usrCfg := params.cfg.UserViperConfig

	id := "latest"
	if params.uuid != "" {
		id = params.uuid
	}

	url := fmt.Sprintf("%s/solutions/%s", usrCfg.GetString("apibaseurl"), id)

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if params.uuid == "" {
		q := req.URL.Query()
		q.Add("exercise_id", params.slug)
		if params.track != "" {
			q.Add("track_id", params.track)
		}
		if params.team != "" {
			q.Add("team_id", params.team)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var payload *downloadPayload
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("unable to parse API response - %s", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		siteURL := config.InferSiteURL(usrCfg.GetString("apibaseurl"))
		return nil, fmt.Errorf("unauthorized request. Please run the configure command. You can find your API token at %s/my/settings", siteURL)
	}

	if res.StatusCode != http.StatusOK {
		switch payload.Error.Type {
		case "track_ambiguous":
			return nil, fmt.Errorf("%s: %s", payload.Error.Message, strings.Join(payload.Error.PossibleTrackIDs, ", "))
		default:
			return nil, errors.New(payload.Error.Message)
		}
	}

	return &downloadPayloadContext{cfg: params.cfg, payload: payload}, nil
}

func (d *downloadPayloadContext) validate() error {
	if d.payload.Error.Message != "" {
		return errors.New(d.payload.Error.Message)
	}
	return nil
}

func (d *downloadPayloadContext) writeMetadata() error {
	if err := d.validate(); err != nil {
		return err
	}

	metadata := d.getMetadata()
	exercise := d.getExercise()

	if err := metadata.Write(exercise.MetadataDir()); err != nil {
		return err
	}

	return nil
}

func (d *downloadPayloadContext) writeSolutionFiles() error {
	if err := d.validate(); err != nil {
		return err
	}
	usrCfg := d.cfg.UserViperConfig

	exercise := d.getExercise()

	for _, file := range d.payload.Solution.Files {
		unparsedURL := fmt.Sprintf("%s%s", d.payload.Solution.FileDownloadBaseURL, file)
		parsedURL, err := netURL.ParseRequestURI(unparsedURL)
		if err != nil {
			return err
		}

		client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
		req, err := client.NewRequest("GET", parsedURL.String(), nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			// TODO: deal with it
			continue
		}
		// Don't bother with empty files.
		if res.Header.Get("Content-Length") == "0" {
			continue
		}

		// TODO: if there's a collision, interactively resolve (show diff, ask if overwrite).
		// TODO: handle --force flag to overwrite without asking.

		// Work around a path bug due to an early design decision (later reversed) to
		// allow numeric suffixes for exercise directories, allowing people to have
		// multiple parallel versions of an exercise.
		pattern := fmt.Sprintf(`\A.*[/\\]%s-\d*/`, exercise.Slug)
		rgxNumericSuffix := regexp.MustCompile(pattern)
		if rgxNumericSuffix.MatchString(file) {
			file = string(rgxNumericSuffix.ReplaceAll([]byte(file), []byte("")))
		}

		// Rewrite paths submitted with an older, buggy client where the Windows path is being treated as part of the filename.
		file = strings.Replace(file, "\\", "/", -1)

		relativePath := filepath.FromSlash(file)

		dir := filepath.Join(exercise.MetadataDir(), filepath.Dir(relativePath))
		if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(exercise.MetadataDir(), relativePath))
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, res.Body); err != nil {
			return err
		}
	}

	return nil
}

func (d *downloadPayloadContext) getMetadata() workspace.ExerciseMetadata {
	return workspace.ExerciseMetadata{
		AutoApprove: d.payload.Solution.Exercise.AutoApprove,
		Track:       d.payload.Solution.Exercise.Track.ID,
		Team:        d.payload.Solution.Team.Slug,
		Exercise:    d.payload.Solution.Exercise.ID,
		ID:          d.payload.Solution.ID,
		URL:         d.payload.Solution.URL,
		Handle:      d.payload.Solution.User.Handle,
		IsRequester: d.payload.Solution.User.IsRequester,
	}
}

func (d *downloadPayloadContext) getExercise() workspace.Exercise {
	usrCfg := d.cfg.UserViperConfig

	root := usrCfg.GetString("workspace")
	if d.payload.Solution.Team.Slug != "" {
		root = filepath.Join(root, "teams", d.payload.Solution.Team.Slug)
	}
	if !d.payload.Solution.User.IsRequester {
		root = filepath.Join(root, "users", d.payload.Solution.User.Handle)
	}

	return workspace.Exercise{
		Root:  root,
		Track: d.payload.Solution.Exercise.Track.ID,
		Slug:  d.payload.Solution.Exercise.ID,
	}
}

type downloadPayload struct {
	Solution struct {
		ID   string `json:"id"`
		URL  string `json:"url"`
		Team struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"team"`
		User struct {
			Handle      string `json:"handle"`
			IsRequester bool   `json:"is_requester"`
		} `json:"user"`
		Exercise struct {
			ID              string `json:"id"`
			InstructionsURL string `json:"instructions_url"`
			AutoApprove     bool   `json:"auto_approve"`
			Track           struct {
				ID       string `json:"id"`
				Language string `json:"language"`
			} `json:"track"`
		} `json:"exercise"`
		FileDownloadBaseURL string   `json:"file_download_base_url"`
		Files               []string `json:"files"`
		Iteration           struct {
			SubmittedAt *string `json:"submitted_at"`
		}
	} `json:"solution"`
	Error struct {
		Type             string   `json:"type"`
		Message          string   `json:"message"`
		PossibleTrackIDs []string `json:"possible_track_ids"`
	} `json:"error,omitempty"`
}
