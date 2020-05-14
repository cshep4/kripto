package model

import "time"

type UserType string

const (
	Landlord UserType = "landlord"
	Tenant   UserType = "tenant"
)

type User struct {
	Id              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	Type            UserType  `json:"type"`
	Verified        bool      `json:"verified"`
	Password        string    `json:"-"`
	Joined          time.Time `json:"-"`
	RentReminders   bool      `json:"-"`
	VerifySignature string    `json:"-"`
	ResetSignature  string    `json:"-"`
}
