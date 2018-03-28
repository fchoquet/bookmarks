package app

import (
	"fmt"
	"strings"
)

// Configuration contains the application Configuration
// We do not consume env vars directly to ease testing
type Configuration struct {
	LogLevel       string
	BasicAuthUsers UserList
	DBConfig       DatabaseConfig
	CSRFSecret     []byte
	// We need this option when using the application over http.
	// Do not enable in prod!!!
	DisableCSRFProtection bool
}

// DatabaseConfig holds the database config and credentials
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

// UserList represents allowed users and passwords
type UserList map[string]string

// ParseUsers the raw env var to get a list of users
// this is a super naive implementation: Semicolumns in passwords are not escaped. Therefore they are not allowed
func ParseUsers(userString string) (UserList, error) {
	users := UserList{}

	if userString == "" {
		return users, nil
	}

	for _, userString := range strings.Split(userString, ";") {
		userInfo := strings.Split(userString, ":")
		if len(userInfo) != 2 {
			return users, fmt.Errorf("invalid basic auth user: %s", userInfo)
		}
		users[strings.Trim(userInfo[0], " ")] = strings.Trim(userInfo[1], " ")
	}
	return users, nil
}
