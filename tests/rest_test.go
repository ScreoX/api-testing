package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRESTPostPatchGet(t *testing.T) {
	client := newRESTClient()
	sandboxQuery := newSandboxQuery(t)

	initial := post{
		UserID: 1,
		Title:  fmt.Sprintf("Go API homework initial title %d", time.Now().UnixNano()),
		Body:   "This body must remain unchanged after PATCH",
	}
	updatedTitle := initial.Title + " - updated"

	createResponse := client.doJSON(t, http.MethodPost, "/posts"+sandboxQuery, initial, nil)
	requireStatus(t, createResponse, http.StatusCreated)
	created := decodeJSON[post](t, createResponse)
	if created.ID == "" {
		t.Fatal("POST /posts returned an empty id")
	}

	resourcePath := "/posts/" + url.PathEscape(created.ID) + sandboxQuery
	patchResponse := client.doJSON(t, http.MethodPatch, resourcePath, map[string]string{"title": updatedTitle}, nil)
	requireStatus(t, patchResponse, http.StatusOK)

	getResponse := client.doJSON(t, http.MethodGet, resourcePath, nil, nil)
	requireStatus(t, getResponse, http.StatusOK)
	actual := decodeJSON[post](t, getResponse)

	if actual.ID != created.ID {
		t.Errorf("GET id=%q, want %q", actual.ID, created.ID)
	}
	if actual.Title != updatedTitle {
		t.Errorf("GET title=%q, want %q", actual.Title, updatedTitle)
	}
	if actual.Body != initial.Body {
		t.Errorf("PATCH changed body: got %q, want %q", actual.Body, initial.Body)
	}
	if actual.UserID != initial.UserID {
		t.Errorf("PATCH changed userId: got %d, want %d", actual.UserID, initial.UserID)
	}
}

func TestBearerAuthentication(t *testing.T) {
	client := newRESTClient()

	username := envOrDefault("API_USERNAME", "Bret")
	login := loginWithCurrentCredentials(t, client, username)
	if strings.TrimSpace(login.AccessToken) == "" {
		t.Fatal("POST /auth/login returned an empty access_token")
	}

	meResponse := client.doJSON(t, http.MethodGet, "/auth/me", nil, map[string]string{
		"Authorization": "Bearer " + login.AccessToken,
	})
	requireStatus(t, meResponse, http.StatusOK)
	me := decodeJSON[user](t, meResponse)

	if me.ID == 0 || !strings.EqualFold(me.Username, username) || me.Email == "" || me.Name == "" {
		t.Fatalf("GET /auth/me returned an unexpected user: %+v", me)
	}
}

func TestRejectsInvalidBearerToken(t *testing.T) {
	client := newRESTClient()

	response := client.doJSON(t, http.MethodGet, "/auth/me", nil, map[string]string{
		"Authorization": "Bearer definitely-not-a-valid-token",
	})
	requireStatus(t, response, http.StatusUnauthorized)

	apiError := decodeJSON[errorResponse](t, response)
	if strings.TrimSpace(apiError.Error) == "" && strings.TrimSpace(apiError.Message) == "" {
		t.Fatalf("invalid-token rejection has no explanatory error: %s", response.body)
	}
}
