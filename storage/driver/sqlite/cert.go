package sqlite

import (
	"fmt"

	"github.com/kassisol/twic/storage/driver"
)

func (c *Config) AddCert(name, ptype, cn, altNames, tsaUrl string) error {
	result := c.DB.Create(&Cert{
		Name:     name,
		Type:     ptype,
		CN:       cn,
		AltNames: altNames,
		TSAURL:   tsaUrl,
	})
	return result.Error
}

func (c *Config) RemoveCert(name string) error {
	if c.certUsedInProfile(name) {
		return fmt.Errorf("cert %q cannot be removed: it is being used by a profile", name)
	}

	result := c.DB.Where("name = ?", name).Delete(&Cert{})
	return result.Error
}

func (c *Config) GetCert(name string) (driver.CertResult, error) {
	var cert Cert

	result := c.DB.Where("name = ?", name).First(&cert)
	if result.Error != nil {
		return driver.CertResult{}, result.Error
	}

	return driver.CertResult{
		Name:     cert.Name,
		Type:     cert.Type,
		CN:       cert.CN,
		AltNames: cert.AltNames,
		TSAURL:   cert.TSAURL,
	}, nil
}

func (c *Config) ListCerts() ([]driver.CertResult, error) {
	var certs []Cert

	result := c.DB.Find(&certs)
	if result.Error != nil {
		return nil, result.Error
	}

	var results []driver.CertResult
	for _, cert := range certs {
		results = append(results, driver.CertResult{
			Name:     cert.Name,
			Type:     cert.Type,
			CN:       cert.CN,
			AltNames: cert.AltNames,
			TSAURL:   cert.TSAURL,
		})
	}

	return results, nil
}

func (c *Config) certUsedInProfile(name string) bool {
	var count int64

	c.DB.Table("profiles").
		Joins("JOIN certs ON certs.id = profiles.cert_id").
		Where("certs.name = ?", name).
		Count(&count)

	return count > 0
}
