package driver

type Storager interface {
	AddCert(name, ptype, cn, altNames, caUrl string) error
	RemoveCert(name string) error
	GetCert(name string) (CertResult, error)
	ListCerts() ([]CertResult, error)

	AddProfile(name, certName, dockerHost string) error
	RemoveProfile(name string) error
	GetProfile(name string) (ProfileResult, error)
	ListProfiles() ([]ProfileResult, error)

	End()
}
