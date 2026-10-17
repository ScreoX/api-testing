package tests

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"testing"
)

const defaultRESTBaseURL = "https://playground.nileslabs.com/api/v1"

type post struct {
	ID     string `json:"id,omitempty"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func newRESTClient() *apiClient {
	return newAPIClient(envOrDefault("REST_BASE_URL", defaultRESTBaseURL))
}

func newSandboxQuery(t *testing.T) string {
	t.Helper()

	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		t.Fatalf("generate sandbox id: %v", err)
	}
	return fmt.Sprintf("?_sandbox=%x", id)
}

func loginWithCurrentCredentials(t *testing.T, client *apiClient, username string) loginResponse {
	t.Helper()

	password := envOrDefault("API_PASSWORD", "Password@123")
	response := client.doJSON(t, http.MethodPost, "/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, nil)
	requireStatus(t, response, http.StatusOK)

	return decodeJSON[loginResponse](t, response)
}
