package version

import (
	"os"
	"runtime"
	"strconv"
	"text/template"
	"time"
)

var (
	Version   string
	GitCommit string
	GitState  string
	BuildDate string
)

var versionTemplate = `Version:     {{.Version}}
Git commit:  {{.GitCommit}}{{if eq .GitState "dirty"}}
Git State:   {{.GitState}}{{end}}
Built:       {{.BuildDate}}
Go version:  {{.GoVersion}}
OS/Arch:     {{.Os}}/{{.Arch}}
`

type VersionInfo struct {
	Version   string
	GoVersion string
	GitCommit string
	GitState  string
	BuildDate string
	Os        string
	Arch      string
}

func New() *VersionInfo {
	buildDate := "unknown"
	if BuildDate != "" {
		if i, err := strconv.ParseInt(BuildDate, 10, 64); err == nil {
			buildDate = time.Unix(i, 0).String()
		}
	}

	version := Version
	if version == "" {
		version = "dev"
	}

	gitCommit := GitCommit
	if gitCommit == "" {
		gitCommit = "unknown"
	}

	return &VersionInfo{
		Version:   version,
		GoVersion: runtime.Version(),
		GitCommit: gitCommit,
		GitState:  GitState,
		BuildDate: buildDate,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func (i *VersionInfo) ShowVersion() {
	tmpl, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, i); err != nil {
		panic(err)
	}
}
