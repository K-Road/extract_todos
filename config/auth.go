package config

import (
	"os"
)

type AuthUser struct {
	Owner    string   //Github user or owner
	APIKey   string   //token or key from secrets
	Projects []string //repos names
}

var Users []AuthUser

func LoadUsersFromEnv() error {
	Users = []AuthUser{
		{
			Owner:    "K-Road",
			APIKey:   os.Getenv("API_KEY"),
			Projects: []string{"*"},
		},
	}
	return nil
}

func GetUsers() []AuthUser {
	return Users
}
