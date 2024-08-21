package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"sync"
	"time"

	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/config"
	"github.com/exercism/cli/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fullAPIKey flag for troubleshoot command.
var fullAPIKey bool

// troubleshootCmd does a diagnostic self-check.
var troubleshootCmd = &cobra.Command{
	Use:     "troubleshoot",
	Aliases: []string{"debug"},
	Short:   "Troubleshoot does a diagnostic self-check.",
	Long: `Provides output to help with troubleshooting.

If you're running into trouble, copy and paste the output from the troubleshoot
command into a topic on the Exercism forum so we can help figure out what's going on.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.TimeoutInSeconds = cli.TimeoutInSeconds * 2
		c := cli.New(Version)

		cfg := config.NewConfig()

		v := viper.New()
		v.AddConfigPath(cfg.Dir)
		v.SetConfigName("user")
		v.SetConfigType("json")
		// Ignore error. If the file doesn't exist, that is fine.
		_ = v.ReadInConfig()

		cfg.UserViperConfig = v

		status := newStatus(c, cfg)
		status.Censor = !fullAPIKey
		s, err := status.check()
		if err != nil {
			return err
		}

		fmt.Printf("%s", s)
		return nil
	},
}

// Status represents the results of a CLI self test.
type Status struct {
	Censor          bool
	Version         versionStatus
	System          systemStatus
	Configuration   configurationStatus
	APIReachability apiReachabilityStatus
	cfg             config.Config
	cli             *cli.CLI
}

type versionStatus struct {
	Current  string
	Latest   string
	Status   string
	Error    error
	UpToDate bool
}

type systemStatus struct {
	OS           string
	Architecture string
	Build        string
}

type configurationStatus struct {
	Home      string
	Workspace string
	Dir       string
	Token     string
	TokenURL  string
}

type apiReachabilityStatus struct {
	Services []*apiPing
}

type apiPing struct {
	Service string
	URL     string
	Status  string
	Latency time.Duration
}

// newStatus prepares a value to perform a diagnostic self-check.
func newStatus(cli *cli.CLI, cfg config.Config) Status {
	status := Status{
		cfg: cfg,
		cli: cli,
	}
	return status
}

// check runs the CLI's diagnostic self-check.
func (status *Status) check() (string, error) {
	status.Version = newVersionStatus(status.cli)
	status.System = newSystemStatus()
	status.Configuration = newConfigurationStatus(status)
	status.APIReachability = newAPIReachabilityStatus(status.cfg)

	return status.compile()
}
func (status *Status) compile() (string, error) {
	t, err := template.New("self-test").Parse(tmplSelfTest)
	if err != nil {
		return "", err
	}

	var bb bytes.Buffer
	if err = t.Execute(&bb, status); err != nil {
		return "", err
	}
	return bb.String(), nil
}

func newAPIReachabilityStatus(cfg config.Config) apiReachabilityStatus {
	baseURL := cfg.UserViperConfig.GetString("apibaseurl")
	if baseURL == "" {
		baseURL = cfg.DefaultBaseURL
	}
	ar := apiReachabilityStatus{
		Services: []*apiPing{
			{Service: "GitHub", URL: "https://api.github.com"},
			{Service: "Exercism", URL: fmt.Sprintf("%s/ping", baseURL)},
		},
	}
	var wg sync.WaitGroup
	wg.Add(len(ar.Services))
	for _, service := range ar.Services {
		go service.Call(&wg)
	}
	wg.Wait()
	return ar
}

func newVersionStatus(c *cli.CLI) versionStatus {
	vs := versionStatus{
		Current: c.Version,
	}
	ok, err := c.IsUpToDate()
	if err == nil {
		vs.Latest = c.LatestRelease.Version()
	} else {
		vs.Error = fmt.Errorf("Error: %s", err)
	}
	vs.UpToDate = ok
	return vs
}

func newSystemStatus() systemStatus {
	ss := systemStatus{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
	if cli.BuildOS != "" && cli.BuildARCH != "" {
		ss.Build = fmt.Sprintf("%s/%s", cli.BuildOS, cli.BuildARCH)
	}
	if cli.BuildARM != "" {
		ss.Build = fmt.Sprintf("%s ARMv%s", ss.Build, cli.BuildARM)
	}
	return ss
}

func newConfigurationStatus(status *Status) configurationStatus {
	v := status.cfg.UserViperConfig

	workspace := v.GetString("workspace")
	if workspace == "" {
		workspace = fmt.Sprintf("%s (default)", config.DefaultWorkspaceDir(status.cfg))
	}

	cs := configurationStatus{
		Home:      status.cfg.Home,
		Workspace: workspace,
		Dir:       status.cfg.Dir,
		Token:     v.GetString("token"),
		TokenURL:  config.SettingsURL(v.GetString("apibaseurl")),
	}
	if status.Censor && cs.Token != "" {
		cs.Token = debug.Redact(cs.Token)
	}
	return cs
}

func (ping *apiPing) Call(wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now()
	res, err := cli.HTTPClient.Get(ping.URL)
	delta := time.Since(now)
	ping.Latency = delta
	if err != nil {
		ping.Status = err.Error()
		return
	}
	res.Body.Close()
	ping.Status = "connected"
}

const tmplSelfTest = `
Troubleshooting Information
===========================

Version
----------------
Current: {{ .Version.Current }}
Latest:  {{ with .Version.Latest }}{{ . }}{{ else }}<unknown>{{ end }}
{{ with .Version.Error }}
{{ . }}
{{ end -}}
{{ if not .Version.UpToDate }}
Call 'exercism upgrade' to get the latest version.
See the release notes at https://github.com/exercism/cli/releases/tag/v{{ .Version.Latest }} for details.
{{ end }}

Operating System
----------------
OS:           {{ .System.OS }}
Architecture: {{ .System.Architecture }}
{{ with .System.Build }}
Build: {{ . }}
{{ end }}

Configuration
----------------
Home:      {{ .Configuration.Home }}
Workspace: {{ .Configuration.Workspace }}
Config:    {{ .Configuration.Dir }}
API key:   {{ with .Configuration.Token }}{{ . }}{{ else }}<not configured>
Find your API key at {{ .Configuration.TokenURL }}{{ end }}

API Reachability
----------------
{{ range .APIReachability.Services }}
{{ .Service }}:
    * {{ .URL }}
    * [{{ .Status }}]
    * {{ .Latency }}
{{ end }}

If you are having trouble, please create a new topic in the Exercism forum
at https://forum.exercism.org/c/support/cli/10 and include
this information.
{{ if not .Censor }}
Don't share your API key. Keep that private.
{{ end }}`

func init() {
	RootCmd.AddCommand(troubleshootCmd)
	troubleshootCmd.Flags().BoolVarP(&fullAPIKey, "full-api-key", "f", false, "display the user's full API key, censored by default")
}
