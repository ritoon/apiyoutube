package middleware

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestGenerateJWT(t *testing.T) {
	data := []struct {
		title       string
		inSecret    string
		inUUID      string
		inExp       time.Time
		expected    *CustomClaims
		expectedErr bool
	}{
		{"valid JWT", "toto", "E5C057FD-6ED8-4A39-9EE2-E18941FCF86F",
			time.Date(2029, time.November, 10, 23, 0, 0, 0, time.UTC),
			&CustomClaims{
				"E5C057FD-6ED8-4A39-9EE2-E18941FCF86F",
				jwt.StandardClaims{
					ExpiresAt: time.Date(2029, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(),
				},
			},
			false,
		},
		{"Expired JWT", "toto", "E5C057FD-6ED8-4A39-9EE2-E18941FCF86F",
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			nil,
			true,
		},
	}
	for _, d := range data {
		got := GenerateJWT(d.inSecret, d.inUUID, d.inExp)
		claim, err := parseJWT(d.inSecret, got)
		if err != nil {
			if d.expectedErr {
				continue
			}
			t.Errorf("for test %v try to get %v and got %v", d.title, d.expected, claim)
			continue
		}

		if claim.UUID != d.expected.UUID || claim.ExpiresAt != d.expected.ExpiresAt {
			t.Errorf("for test %v try to get %v and got %v", d.title, d.expected, claim)
		}
	}
}
