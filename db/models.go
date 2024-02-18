// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ValidRoles string

const (
	ValidRolesAdmin ValidRoles = "Admin"
	ValidRolesUser  ValidRoles = "User"
)

func (e *ValidRoles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ValidRoles(s)
	case string:
		*e = ValidRoles(s)
	default:
		return fmt.Errorf("unsupported scan type for ValidRoles: %T", src)
	}
	return nil
}

type NullValidRoles struct {
	ValidRoles ValidRoles
	Valid      bool // Valid is true if ValidRoles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullValidRoles) Scan(value interface{}) error {
	if value == nil {
		ns.ValidRoles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ValidRoles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullValidRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ValidRoles), nil
}

type Bucket struct {
	BucketID       int32
	Name           string
	UserID         pgtype.Int4
	OrganizationID pgtype.Int4
	Envvariables   []byte
}

type Organization struct {
	OrganizationID int32
	Name           string
}

type User struct {
	UserID         int32
	OrganizationID pgtype.Int4
	Email          string
	Password       string
	Role           NullValidRoles
}
