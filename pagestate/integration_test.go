package pagestate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"page-state-saver/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagestateIntegration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)

	assert.NoError(t, err)

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("error terminating postgres container: %s", err)
		}
	}()

	repository, err := NewRepository(ctx, pgContainer.ConnectionString)

	assert.NoError(t, err)

	server := NewServer("8080", repository)

	ts := httptest.NewServer(http.HandlerFunc(server.server.Handler.ServeHTTP))

	defer ts.Close()

	requestBody := map[string]interface{}{
		"url":       "https://example.com",
		"scrollPos": 33,
	}

	jsonData, err := json.Marshal(requestBody)

	assert.NoError(t, err)

	resp, err := http.Post(ts.URL+"/pagestate", "application/json", bytes.NewBuffer(jsonData))

	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response CreatePagestateResponse

	err = json.NewDecoder(resp.Body).Decode(&response)

	assert.NoError(t, err)
	assert.NotZero(t, response.Id)
	assert.Equal(t, "https://example.com", response.Url)
	assert.Equal(t, 33, response.ScrollPos)
}
