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

func (suite *PagestateAPITestSuite) TestSavePagestateAPI() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	requestBody := map[string]any{
		"url":         "https://example.com",
		"scrollPos":   33,
		"visibleText": "Sample visible text content",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	resp, err := http.Post(ts.URL+"/pagestate/save", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response CreatePagestateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.NotZero(t, response.Id)
	assert.Equal(t, "https://example.com", response.Url)
	assert.Equal(t, 33, response.ScrollPos)
	assert.Equal(t, "Sample visible text content", response.VisibleText)
}

func TestPagestateAPITestSuite(t *testing.T) {
	suite.Run(t, new(PagestateAPITestSuite))
}
