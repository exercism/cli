package cli

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/exercism/cli/user"
)

// Status represents the results of a CLI self test.
type Status struct {
	Censor          bool
	Version         versionStatus
	System          systemStatus
	Configuration   configurationStatus
	APIReachability apiReachabilityStatus
	cli             *CLI
	cfg             user.Config
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
	File      string
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

func NewStatus(c *CLI, uc user.Config) Status {
	status := Status{
		cli: c,
		cfg: uc,
	}
	return status
}

func (status *Status) Check() (string, error) {
	status.Version = newVersionStatus(status.cli)
	status.System = newSystemStatus()
	status.Configuration = newConfigurationStatus(status)
	status.APIReachability = newAPIReachabilityStatus()

	return status.compile()
}
func (status *Status) compile() (string, error) {
	t, err := template.New("self-test").Parse(tmplSelfTest)
	if err != nil {
		return "", err
	}

	var bb bytes.Buffer
	t.Execute(&bb, status)
	return bb.String(), nil
}

func newAPIReachabilityStatus() apiReachabilityStatus {
	ar := apiReachabilityStatus{
		Services: []*apiPing{
			{Service: "GitHub", URL: "https://api.github.com"},
			{Service: "Exercism", URL: "http://exercism.io/api/v1"},
			{Service: "X-API", URL: "http://x.exercism.io"},
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

func newVersionStatus(cli *CLI) versionStatus {
	vs := versionStatus{
		Current: cli.Version,
	}
	ok, err := cli.IsUpToDate()
	if err == nil {
		vs.Latest = cli.LatestRelease.Version()
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
	if BuildOS != "" && BuildARCH != "" {
		ss.Build = fmt.Sprintf("%s/%s", BuildOS, BuildARCH)
	}
	if BuildARM != "" {
		ss.Build = fmt.Sprintf("%s ARMv%s", ss.Build, BuildARM)
	}
	return ss
}

func newConfigurationStatus(status *Status) configurationStatus {
	cs := configurationStatus{
		Home:      status.cfg.Home,
		Workspace: status.cfg.Workspace,
		File:      status.cfg.Path,
		Token:     status.cfg.Token,
		TokenURL:  "http://exercism.io/account/key",
	}
	if status.Censor {
		cs.Token = redactToken(status.cfg.Token)
	}
	return cs
}

func (ping *apiPing) Call(wg *sync.WaitGroup) {
	defer wg.Done()

	now := time.Now()
	res, err := HTTPClient.Get(ping.URL)
	delta := time.Since(now)
	ping.Latency = delta
	if err != nil {
		ping.Status = err.Error()
		return
	}
	res.Body.Close()
	ping.Status = "connected"
}

func redactToken(token string) string {
	str := token[4 : len(token)-3]
	redaction := strings.Repeat("*", len(str))
	return string(token[:4]) + redaction + string(token[len(token)-3:])
}

const tmplSelfTest = `

Debug Information
=================

Version
----------------
Current: {{ .Version.Current }}
Latest:  {{ with .Version.Latest }}{{ . }}{{ else }}<unknown>{{ end }}
{{ with .Version.Error }}
{{ . }}
{{ end -}}
{{ if not .Version.UpToDate }}
Call 'exercism upgrade' to get the latest version.
See the release notes at https://github.com/exercism/cli/releases/tag/{{ .Version.Latest }} for details.
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
Config:    {{ .Configuration.File }}
API key:   {{ with .Configuration.Token }}{{ . }}{{ else }}<not configured>
Find your API key at {{ .Configuration.TokenURL }}{{ end }}

API Reachability
----------------
{{ range .APIReachability.Services }}
{{ .Service }}:
          {{ .URL }}
          [{{ .Status }}]
          {{ .Latency }}
{{ end }}

If you are having trouble please file a GitHub issue at
https://github.com/exercism/exercism.io/issues and include
this information.
{{ if not .Censor }}
Don't share your API key. Keep that private.
{{ end }}
`
