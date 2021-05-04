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

func (ns MyNullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

func (ns *MyNullString) UnmarshalJSON(b []byte) error {
	if "null" == string(b) {
		ns.Valid = false
	} else {

		err := json.Unmarshal(b, &ns.String)

		if err != nil {
			return err
		} else {
			ns.Valid = true
		}
	}

	return nil

}

func (ns *MyNullString) Value() interface{} {

	if ns.Valid {
		return ns.String
	} else {
		return ns.NullString
	}
}

type MyNullTime struct {
	sql.NullTime
}

func (n MyNullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return []byte(`null`), nil

}

func (ns *MyNullTime) UnmarshalJSON(b []byte) error {

	if "null" == string(b) {
		ns.Valid = false
	} else {

		err := json.Unmarshal(b, &ns.Time)

		if err != nil {
			ns.Valid = false
			return err
		} else {
			ns.Valid = true
		}

	}

	return nil

}

func (ns *MyNullTime) Value() interface{} {

	if ns.Valid {
		return ns.Time
	} else {
		return ns.NullTime
	}

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
