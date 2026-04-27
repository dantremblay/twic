package profile

import (
	"errors"
	"os"
	"text/template"

	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

type Data struct {
	Shell     string
	Unset     bool
	TLSVerify string
	Host      string
	CertPath  string
}

func newEnvCommand() *cobra.Command {
	var (
		shell string
		unset bool
	)

	cmd := &cobra.Command{
		Use:   "env [name]",
		Short: "Set / Unset Docker environment variables",
		Long:  envDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnv(args, shell, unset)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&shell, "shell", "s", "bash", "Force environment to be configured for a specified shell: (tcsh, bash)")
	flags.BoolVarP(&unset, "unset", "u", false, "Unset variables instead of setting them")

	return cmd
}

func runEnv(args []string, shell string, unset bool) error {
	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
	}

	if shell != "bash" && shell != "tcsh" {
		return errors.New("shell is not correct")
	}

	cfg := adf.NewClient()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	profile, err := s.GetProfile(args[0])
	if err != nil {
		return errors.New("profile does not exist")
	}

	cfg.SetName(profile.Cert.Name)

	data := Data{
		Shell:     shell,
		Unset:     unset,
		TLSVerify: "1",
		Host:      profile.DockerHost,
		CertPath:  cfg.Profile.CertDir,
	}

	t := template.Must(template.New("Shell commands template").Parse(envTpl))
	return t.Execute(os.Stdout, data)
}

var envTpl = `
{{- if eq .Shell "bash" }}
{{- if .Unset }}
unset DOCKER_HOST DOCKER_TLS_VERIFY DOCKER_CERT_PATH
{{- else }}
export DOCKER_HOST={{ .Host }}
export DOCKER_TLS_VERIFY={{ .TLSVerify }}
export DOCKER_CERT_PATH={{ .CertPath }}/
{{- end }}
{{- end }}
{{- if eq .Shell "tcsh" }}
{{- if .Unset }}
unsetenv DOCKER_HOST DOCKER_TLS_VERIFY DOCKER_CERT_PATH
{{- else }}
setenv DOCKER_HOST {{ .Host }}
setenv DOCKER_TLS_VERIFY {{ .TLSVerify }}
setenv DOCKER_CERT_PATH {{ .CertPath }}/
{{- end }}
{{- end }}
`

var envDescription = `
Set / Unset Docker environment variables

`
