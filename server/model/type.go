package model

import (
	"database/sql"
	"encoding/json"

	"github.com/Netflix/go-env"
)

/**
 */
type MyNullString struct {
	sql.NullString
}
type MyNullTime struct {
	sql.NullTime
}

func (ns MyNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

func (n MyNullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return []byte(`null`), nil

}

type Environment struct {
	Debug        bool   `env:"TRACKING_DATABASE_DEBUG"`
	Password     string `env:"TRACKING_DATABASE_PASSWORD"`
	DatabasePort int    `env:"TRACKING_DATABASE_PORT"`
	Server       string `env:"TRACKING_DATABASE_ADDRESS"`
	User         string `env:"TRACKING_DATABASE_USER"`
	Database     string `env:"TRACKING_DATABASE_NAME"`
	HttpPort     int    `env:"ASPNETCORE_PORT"`
	Extras       env.EnvSet
}
