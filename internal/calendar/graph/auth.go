package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// clientCredentialsToken fetches an app-level token for OneDrive/Excel
func clientCredentialsToken() (string, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	endpoint := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/token",
		tenantID,
	)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", os.Getenv("AZURE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("AZURE_CLIENT_SECRET"))
	data.Set("scope", "https://graph.microsoft.com/.default")

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("token decode failed: %w", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("token error: %s — %s", result.Error, result.Description)
	}
	return result.AccessToken, nil
}

// GetToken returns an app-level token — used for OneDrive and Excel calls
func GetToken() (string, error) {
	return clientCredentialsToken()
}

// GetCalendarToken returns a delegated token — used for Outlook calendar calls
// It uses the refresh token stored in .env to silently get a fresh access token
func GetCalendarToken() (string, error) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		Scopes:       []string{"Calendars.ReadWrite", "offline_access"},
		Endpoint:     microsoft.AzureADEndpoint("common"),
		RedirectURL:  "http://localhost:8080/callback",
	}

	refreshToken := os.Getenv("OUTLOOK_REFRESH_TOKEN")
	if refreshToken == "" {
		return "", fmt.Errorf("OUTLOOK_REFRESH_TOKEN is not set in .env — run cmd/gettoken first")
	}

	existing := &oauth2.Token{RefreshToken: refreshToken}
	tokenSource := conf.TokenSource(context.Background(), existing)

	newToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("calendar token refresh failed: %w", err)
	}
	return newToken.AccessToken, nil
}

// GraphRequest makes an authenticated HTTP call to Microsoft Graph
func GraphRequest(method, endpoint, token string, body interface{}) (*http.Response, error) {
	var reqBody *strings.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = strings.NewReader(string(b))
	} else {
		reqBody = strings.NewReader("")
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		"https://graph.microsoft.com/v1.0"+endpoint,
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
