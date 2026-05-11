package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
}

type Cert struct {
	Model

	Name     string `gorm:"uniqueIndex"`
	Type     string
	CN       string
	AltNames string

	TSAURL string
}

type Profile struct {
	Model

	Name       string `gorm:"uniqueIndex"`
	Cert       Cert
	CertID     uint
	DockerHost string
}
