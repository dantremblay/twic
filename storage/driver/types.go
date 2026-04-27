package driver

type CertResult struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	CN       string `json:"cn"`
	AltNames string `json:"alt_names,omitempty"`
	TSAURL   string `json:"tsa_url"`
}

type ProfileResult struct {
	Name       string     `json:"name"`
	Cert       CertResult `json:"cert"`
	DockerHost string     `json:"docker_host"`
}
