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
		"url":         "https://example.com",
		"scrollPos":   33,
		"visibleText": "Sample visible text content",
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
	assert.Equal(t, "Sample visible text content", response.VisibleText)
}

func (suite *PagestateAPITestSuite) TestGetPagestateAPI() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	createJsonData1, err := json.Marshal(map[string]any{"url": "https://test1.com", "scrollPos": 150, "visibleText": "First page content"})
	assert.NoError(t, err)

	createResp1, err := http.Post(ts.URL+"/pagestate", "application/json", bytes.NewBuffer(createJsonData1))
	assert.NoError(t, err)
	defer createResp1.Body.Close()

	assert.Equal(t, http.StatusCreated, createResp1.StatusCode)

	var createResponse1 CreatePagestateResponse
	err = json.NewDecoder(createResp1.Body).Decode(&createResponse1)
	assert.NoError(t, err)

	createJsonData2, err := json.Marshal(map[string]any{"url": "https://test2.com", "scrollPos": 300, "visibleText": "Second page content"})
	assert.NoError(t, err)

	createResp2, err := http.Post(ts.URL+"/pagestate", "application/json", bytes.NewBuffer(createJsonData2))
	assert.NoError(t, err)
	defer createResp2.Body.Close()

	assert.Equal(t, http.StatusCreated, createResp2.StatusCode)

	var createResponse2 CreatePagestateResponse
	err = json.NewDecoder(createResp2.Body).Decode(&createResponse2)
	assert.NoError(t, err)

	getResp, err := http.Get(ts.URL + "/pagestate/get")
	assert.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var getResponse GetAllPagestatesResponse
	err = json.NewDecoder(getResp.Body).Decode(&getResponse)
	assert.NoError(t, err)

	assert.Len(t, getResponse.Pagestates, 2)

	found1, found2 := false, false
	for _, retrieved := range getResponse.Pagestates {
		if retrieved.Id == createResponse1.Id {
			assert.Equal(t, createResponse1.Url, retrieved.Url)
			assert.Equal(t, createResponse1.ScrollPos, retrieved.ScrollPos)
			assert.Equal(t, createResponse1.VisibleText, retrieved.VisibleText)
			found1 = true
		}
		if retrieved.Id == createResponse2.Id {
			assert.Equal(t, createResponse2.Url, retrieved.Url)
			assert.Equal(t, createResponse2.ScrollPos, retrieved.ScrollPos)
			assert.Equal(t, createResponse2.VisibleText, retrieved.VisibleText)
			found2 = true
		}
	}
	assert.True(t, found1, "First created page state not found in response")
	assert.True(t, found2, "Second created page state not found in response")
}

func TestPagestateAPITestSuite(t *testing.T) {
	suite.Run(t, new(PagestateAPITestSuite))
}
