package sqlite

import (
	"github.com/kassisol/twic/storage/driver"
)

func (c *Config) AddProfile(name, certName, dockerHost string) error {
	var cert Cert
	if result := c.DB.Where("name = ?", certName).First(&cert); result.Error != nil {
		return result.Error
	}

	result := c.DB.Create(&Profile{
		Name:       name,
		Cert:       cert,
		DockerHost: dockerHost,
	})
	return result.Error
}

func (c *Config) RemoveProfile(name string) error {
	result := c.DB.Where("name = ?", name).Delete(&Profile{})
	return result.Error
}

func (c *Config) GetProfile(name string) (driver.ProfileResult, error) {
	var profile Profile

	result := c.DB.Preload("Cert").Where("name = ?", name).First(&profile)
	if result.Error != nil {
		return driver.ProfileResult{}, result.Error
	}

	return driver.ProfileResult{
		Name: profile.Name,
		Cert: driver.CertResult{
			Name:     profile.Cert.Name,
			Type:     profile.Cert.Type,
			CN:       profile.Cert.CN,
			AltNames: profile.Cert.AltNames,
			TSAURL:   profile.Cert.TSAURL,
		},
		DockerHost: profile.DockerHost,
	}, nil
}

func (c *Config) ListProfiles() ([]driver.ProfileResult, error) {
	var profiles []Profile

	result := c.DB.Preload("Cert").Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}

	var results []driver.ProfileResult
	for _, profile := range profiles {
		results = append(results, driver.ProfileResult{
			Name: profile.Name,
			Cert: driver.CertResult{
				Name:     profile.Cert.Name,
				Type:     profile.Cert.Type,
				CN:       profile.Cert.CN,
				AltNames: profile.Cert.AltNames,
				TSAURL:   profile.Cert.TSAURL,
			},
			DockerHost: profile.DockerHost,
		})
	}

	return results, nil
}
