package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	netURL "net/url"
	"os"
	"path/filepath"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"r"},
	Short:   "Restart an exercise.",
	Long: `Restart an exercise.

You may restart an exercise to work on. If you've already
started working on it, the command will override your local solution.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()
		cfg.UserViperConfig = v

		return runRestart(cfg, cmd.Flags(), args)
	},
}

func runRestart(cfg config.Config, flags *pflag.FlagSet, args []string) error {
	usrCfg := cfg.UserViperConfig
	if err := validateUserConfig(usrCfg); err != nil {
		return err
	}

	download, err := newRestart(flags, usrCfg)
	if err != nil {
		return err
	}

	metadata := download.payload.metadata()
	dir := metadata.Exercise(usrCfg.GetString("workspace")).MetadataDir()

	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return err
	}

	if err := metadata.Write(dir); err != nil {
		return err
	}

	client, err := api.NewClient(usrCfg.GetString("token"), usrCfg.GetString("apibaseurl"))
	if err != nil {
		return err
	}

	for _, sf := range download.payload.files() {
		url, err := sf.url()
		if err != nil {
			return err
		}

		req, err := client.NewRequest("GET", url, nil)
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

		path := sf.relativePath()
		dir := filepath.Join(metadata.Dir, filepath.Dir(path))
		if err = os.MkdirAll(dir, os.FileMode(0755)); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(metadata.Dir, path))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, res.Body)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(Err, "\nDownloaded to\n")
	fmt.Fprintf(Out, "%s\n", metadata.Dir)
	return nil
}

type restart struct {
	// exercise
	slug string

	// user config
	token, apibaseurl, workspace string

	// optional
	track string

	payload *downloadPayload
}

func newRestart(flags *pflag.FlagSet, usrCfg *viper.Viper) (*restart, error) {
	var err error
	r := &restart{}
	r.slug, err = flags.GetString("exercise")
	if err != nil {
		return nil, err
	}
	r.track, err = flags.GetString("track")
	if err != nil {
		return nil, err
	}

	r.token = usrCfg.GetString("token")
	r.apibaseurl = usrCfg.GetString("apibaseurl")
	r.workspace = usrCfg.GetString("workspace")

	if err = r.needsSlug(); err != nil {
		return nil, err
	}
	if err = r.needsUserConfigValues(); err != nil {
		return nil, err
	}
	if err = r.needsSlugWhenGivenTrack(); err != nil {
		return nil, err
	}

	client, err := api.NewClient(r.token, r.apibaseurl)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", r.url(), nil)
	if err != nil {
		return nil, err
	}
	r.buildQueryParams(req.URL)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, decodedAPIError(res)
	}

	body, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(body))

	if err := json.Unmarshal(body, &r.payload); err != nil {
		return nil, decodedAPIError(res)
	}

	return r, nil
}

func (r restart) url() string {
	id := "latest"
	return fmt.Sprintf("%s/solutions/%s", r.apibaseurl, id)
}

func (r restart) buildQueryParams(url *netURL.URL) {
	query := url.Query()
	if r.slug != "" {
		query.Add("exercise_id", r.slug)
		if r.track != "" {
			query.Add("track_id", r.track)
		}
	}
	url.RawQuery = query.Encode()
}

// needsSlug checks the presence of slug.
func (r restart) needsSlug() error {
	if r.slug != "" {
		return errors.New("need an --exercise name")
	}
	return nil
}

// needsUserConfigValues checks the presence of required values from the user config.
func (r restart) needsUserConfigValues() error {
	errMsg := "missing required user config: '%s'"
	if r.token == "" {
		return fmt.Errorf(errMsg, "token")
	}
	if r.apibaseurl == "" {
		return fmt.Errorf(errMsg, "apibaseurl")
	}
	if r.workspace == "" {
		return fmt.Errorf(errMsg, "workspace")
	}
	return nil
}

// needsSlugWhenGivenTrack ensures that track arguments are also given with a slug.
func (r restart) needsSlugWhenGivenTrack() error {
	if (r.track != "") && r.slug == "" {
		return errors.New("--track requires --exercise")
	}
	return nil
}

func setupRestartFlags(flags *pflag.FlagSet) {
	flags.StringP("track", "t", "", "the track ID")
	flags.StringP("exercise", "e", "", "the exercise slug")
}

func init() {
	RootCmd.AddCommand(restartCmd)
	setupRestartFlags(restartCmd.Flags())
}
