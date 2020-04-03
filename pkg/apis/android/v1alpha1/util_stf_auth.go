package v1alpha1

import (
	"fmt"
)

// UseMockAuth returns true if this STF deployment should use the mock
// authentication service. Otherwise we use OAuth.
func (s *STFConfig) UseMockAuth() bool {
	if s.Auth != nil {
		if s.Auth.Mock {
			return true
		} else if s.Auth.OAuth != nil {
			return false
		}
	}
	return true
}

// GetAuthURL returns what the authentication URL should be for this STF
// deployment.
func (s *STFConfig) GetAuthURL() string {
	var authType string
	if s.UseMockAuth() {
		authType = "mock"
	} else {
		authType = "oauth"
	}
	return fmt.Sprintf("%s/auth/%s/", s.GetAppExternalURL(), authType)
}

// GetOAuthClientIDKey returns where in our STFSecret we can find our client ID
// for OAuth.
func (s *STFConfig) GetOAuthClientIDKey() string {
	if s.Auth != nil {
		if s.Auth.OAuth != nil {
			if s.Auth.OAuth.ClientIDKey != "" {
				return s.Auth.OAuth.ClientIDKey
			}
		}
	}
	return defaultOAuthClientIDKey
}

// GetOAuthClientSecretKey returns where in our STFSecret we can find the client
// secret for OAuth.
func (s *STFConfig) GetOAuthClientSecretKey() string {
	if s.Auth != nil {
		if s.Auth.OAuth != nil {
			if s.Auth.OAuth.ClientSecretKey != "" {
				return s.Auth.OAuth.ClientSecretKey
			}
		}
	}
	return defaultOAuthClientSecretKey
}
