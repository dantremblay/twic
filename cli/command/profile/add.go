package profile

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/pkg/urlutil"
	"github.com/kassisol/twic/pkg/validate"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	var (
		certName     string
		dockerScheme string
		dockerHost   string
		dockerPort   string
	)

	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add Docker profile",
		Long:  addDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args, certName, dockerScheme, dockerHost, dockerPort)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&certName, "cert-name", "c", "", "Certificate Name")
	flags.StringVarP(&dockerScheme, "docker-scheme", "s", "tcp", "Docker Scheme")
	flags.StringVarP(&dockerHost, "docker-host", "a", "", "Docker Host")
	flags.StringVarP(&dockerPort, "docker-port", "p", "2376", "Docker Port")

	return cmd
}

func runAdd(args []string, certName, dockerScheme, dockerHost, dockerPort string) error {
	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
	}

	certname := certName
	if len(certname) == 0 {
		certname = input.ReadInput("Certificate Name")
	}

	dockerscheme := dockerScheme
	if len(dockerscheme) == 0 {
		dockerscheme = input.ReadInput("Docker Scheme")
	}

	dockerhost := dockerHost
	if len(dockerhost) == 0 {
		dockerhost = input.ReadInput("Docker Host")
	}

	dockerport := dockerPort
	if len(dockerport) == 0 {
		dockerport = input.ReadInput("Docker Port")
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

	if err := validate.Name(args[0]); err != nil {
		return err
	}

	cert, err := s.GetCert(certname)
	if err != nil || len(cert.Name) == 0 {
		return errors.New("certificate name is not valid")
	}

	if cert.Type == "engine" {
		return errors.New("engine certificate type cannot be added to profile")
	}

	if dockerscheme != "tcp" {
		return errors.New("docker host scheme should be tcp")
	}

	if len(dockerhost) == 0 {
		return errors.New("docker host cannot be empty")
	}

	p, err := strconv.Atoi(dockerport)
	if err != nil {
		return err
	}
	if err := validate.Port(p); err != nil {
		return err
	}

	dockerurl := fmt.Sprintf("%s://%s:%s", dockerscheme, dockerhost, dockerport)
	if _, err := urlutil.Parse(dockerurl); err != nil {
		return err
	}

	return s.AddProfile(args[0], certname, dockerurl)
}

var addDescription = `
Add Docker profile

`
