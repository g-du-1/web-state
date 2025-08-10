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
	_, err := suite.repository.pool.Exec(suite.ctx, "TRUNCATE TABLE pagestates RESTART IDENTITY CASCADE")

	if err != nil {
		suite.T().Logf("error clearing database: %s", err)
	}
}

func (suite *PagestateAPITestSuite) TearDownSuite() {
	suite.repository.Close()

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

	resp, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response PagestateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.NotZero(t, response.Id)
	assert.Equal(t, "https://example.com", response.Url)
	assert.Equal(t, 33, response.ScrollPos)
	assert.Equal(t, "Sample visible text content", response.VisibleText)
}

func (suite *PagestateAPITestSuite) TestGetsExistingPageStateForUrl() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	url := "https://example.com/unique"

	requestBody1 := map[string]any{
		"url":         url,
		"scrollPos":   1,
		"visibleText": "first",
	}

	jsonData1, _ := json.Marshal(requestBody1)
	resp1, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData1))
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	requestBody2 := map[string]any{
		"url":         url,
		"scrollPos":   2,
		"visibleText": "second",
	}

	jsonData2, _ := json.Marshal(requestBody2)
	resp2, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData2))
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)

	getResp, err := http.Get(ts.URL + "/api/v1/pagestate/all")
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var response []PagestateResponse
	err = json.NewDecoder(getResp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Len(t, response, 1)
	assert.Equal(t, url, response[0].Url)
	assert.Equal(t, 2, response[0].ScrollPos)
	assert.Equal(t, "second", response[0].VisibleText)
}

func (suite *PagestateAPITestSuite) TestGetsAllPageStatesInOrder() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	url1 := "https://example.com/some-page"
	url2 := "https://example.com/some-other-page"

	requestBody1 := map[string]any{
		"url":         url1,
		"scrollPos":   9876,
		"visibleText": "some text",
	}

	requestBody2 := map[string]any{
		"url":         url2,
		"scrollPos":   4321,
		"visibleText": "other text",
	}

	jsonData1, err := json.Marshal(requestBody1)
	assert.NoError(t, err)
	resp1, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData1))
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	jsonData2, err := json.Marshal(requestBody2)
	assert.NoError(t, err)
	resp2, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData2))
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)

	getResp, err := http.Get(ts.URL + "/api/v1/pagestate/all")
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var response []PagestateResponse
	err = json.NewDecoder(getResp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	assert.Equal(t, url2, response[0].Url)
	assert.Equal(t, 4321, response[0].ScrollPos)
	assert.Equal(t, "other text", response[0].VisibleText)

	assert.Equal(t, url1, response[1].Url)
	assert.Equal(t, 9876, response[1].ScrollPos)
	assert.Equal(t, "some text", response[1].VisibleText)
}

func (suite *PagestateAPITestSuite) TestDeletesAllPageStates() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	url1 := "https://example.com/page1"
	url2 := "https://example.com/page2"

	body1, err := json.Marshal(map[string]any{
		"url":         url1,
		"scrollPos":   1,
		"visibleText": "a",
	})

	assert.NoError(t, err)

	body2, err := json.Marshal(map[string]any{
		"url":         url2,
		"scrollPos":   2,
		"visibleText": "b",
	})

	assert.NoError(t, err)

	resp1, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(body1))
	assert.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	resp2, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(body2))
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/pagestate/delete", nil)
	assert.NoError(t, err)
	deleteResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer deleteResp.Body.Close()

	getResp, err := http.Get(ts.URL + "/api/v1/pagestate/all")
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var response []PagestateResponse
	err = json.NewDecoder(getResp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response)
}

func (suite *PagestateAPITestSuite) TestReturns404WhenPageStateDoesNotExist() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	getResp, err := http.Get(ts.URL + "/api/v1/pagestate?url=https://example.com/page1")
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

func (suite *PagestateAPITestSuite) TestReturnsPageStateWhenExists() {
	t := suite.T()

	ts := httptest.NewServer(http.HandlerFunc(suite.server.server.Handler.ServeHTTP))
	defer ts.Close()

	url := "https://example.com/test-page"

	requestBody := map[string]any{
		"url":         url,
		"scrollPos":   100,
		"visibleText": "Test content",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	saveResp, err := http.Post(ts.URL+"/api/v1/pagestate/save", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer saveResp.Body.Close()
	assert.Equal(t, http.StatusCreated, saveResp.StatusCode)

	getResp, err := http.Get(ts.URL + "/api/v1/pagestate?url=" + url)
	assert.NoError(t, err)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var response PagestateResponse
	err = json.NewDecoder(getResp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, url, response.Url)
	assert.Equal(t, 100, response.ScrollPos)
	assert.Equal(t, "Test content", response.VisibleText)
}

func TestPagestateAPITestSuite(t *testing.T) {
	suite.Run(t, new(PagestateAPITestSuite))
}
