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
	"github.com/stretchr/testify/suite"
)

type PagestateAPITestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *Repository
	server      *Server
	ctx         context.Context
}

func (suite *PagestateAPITestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)

	if err != nil {
		suite.T().Fatal(err)
	}

	suite.pgContainer = pgContainer
	repository, err := NewRepository(suite.ctx, suite.pgContainer.ConnectionString)

	if err != nil {
		suite.T().Fatal(err)
	}

	suite.repository = repository
	suite.server = NewServer("8080", repository)
}

func (suite *PagestateAPITestSuite) SetupTest() {
	_, err := suite.repository.conn.Exec(suite.ctx, "TRUNCATE TABLE pagestates RESTART IDENTITY CASCADE")

	if err != nil {
		suite.T().Logf("error clearing database: %s", err)
	}
}

func (suite *PagestateAPITestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Logf("error terminating postgres container: %s", err)
	}
}

func (suite *PagestateAPITestSuite) TestCreatePagestateAPI() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
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

func (suite *PagestateAPITestSuite) TestGetPagestateAPI() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	createRequests := []map[string]any{
		{"url": "https://test1.com", "scrollPos": 150},
		{"url": "https://test2.com", "scrollPos": 300},
		{"url": "https://test3.com", "scrollPos": 450},
	}

	var createdPagestates []CreatePagestateResponse

	for _, requestBody := range createRequests {
		createJsonData, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		createResp, err := http.Post(ts.URL+"/pagestate", "application/json", bytes.NewBuffer(createJsonData))
		assert.NoError(t, err)
		defer createResp.Body.Close()

		assert.Equal(t, http.StatusCreated, createResp.StatusCode)

		var createResponse CreatePagestateResponse
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		assert.NoError(t, err)
		createdPagestates = append(createdPagestates, createResponse)
	}

	getResp, err := http.Get(ts.URL + "/pagestate/get")
	assert.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var getResponse GetAllPagestatesResponse
	err = json.NewDecoder(getResp.Body).Decode(&getResponse)
	assert.NoError(t, err)

	assert.Len(t, getResponse.Pagestates, len(createRequests))

	for _, created := range createdPagestates {
		found := false
		for _, retrieved := range getResponse.Pagestates {
			if retrieved.Id == created.Id {
				assert.Equal(t, created.Url, retrieved.Url)
				assert.Equal(t, created.ScrollPos, retrieved.ScrollPos)
				found = true
				break
			}
		}
		assert.True(t, found, "Created page state with ID %d not found in response", created.Id)
	}
}

func TestPagestateAPITestSuite(t *testing.T) {
	suite.Run(t, new(PagestateAPITestSuite))
}
