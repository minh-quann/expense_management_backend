package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// GoogleTokenInfo represents the response from Google tokeninfo endpoint
type GoogleTokenInfo struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"` // Google ID
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Exp           string `json:"exp"`
}

// VerifyGoogleToken verifies a Google ID Token using the tokeninfo API
func VerifyGoogleToken(idToken string) (*GoogleTokenInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid google token")
	}

	var info GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Email == "" {
		return nil, errors.New("email not found in google token")
	}

	return &info, nil
}
