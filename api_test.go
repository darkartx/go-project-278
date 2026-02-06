package code

import (
	"bytes"
	"code/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestLinksList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	req, _ := http.NewRequest("GET", "http://localhost/api/links", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedLinks := []handlers.Link{
		{Id: 1, OriginalUrl: "http://google.com", ShortName: "test", ShortUrl: "http://localhost/r/test"},
		{Id: 2, OriginalUrl: "http://google.com", ShortName: "test2", ShortUrl: "http://localhost/r/test2"},
	}
	var actualLinks []handlers.Link
	err := json.Unmarshal(w.Body.Bytes(), &actualLinks)
	assert.NoError(t, err)
	assert.Equal(t, expectedLinks, actualLinks)
}

func TestLinksCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com","short_name":"testtest"}`
	req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expected := `{"id":1,"original_url":"http://google.com","short_name":"testtest","short_url":"http://localhost/r/testtest"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksCreateWithoutShortName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com"}`
	req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var link handlers.Link
	err := json.Unmarshal(w.Body.Bytes(), &link)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), link.Id)
	assert.Equal(t, "http://google.com", link.OriginalUrl)
	assert.LessOrEqual(t, len(link.ShortName), 10)
	assert.GreaterOrEqual(t, len(link.ShortName), 6)
	assert.Equal(t, "http://localhost/r/"+link.ShortName, link.ShortUrl)
}

func TestLinksCreateWithInvalidOriginalUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"invalid-url","short_name":"testtest"}`
	req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid original url"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksCreateWithInvalidShortName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com","short_name":"!@#$!asdasd"}`
	req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid short name"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	req, _ := http.NewRequest("GET", "http://localhost/api/links/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"id":1,"original_url":"https://google.com","short_name":"test","short_url":"http://localhost/r/test"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksGetWithInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	req, _ := http.NewRequest("GET", "http://localhost/api/links/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"id":1,"original_url":"http://google.com","short_name":"testtest","short_url":"http://localhost/r/testtest"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/abc", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithInvalidOriginalUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"invalid-url","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid original url"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithInvalidShortName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	body := `{"original_url":"http://google.com","short_name":"!@#$!asdasd"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid short name"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	req, _ := http.NewRequest("DELETE", "http://localhost/api/links/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestLinksDeleteWithInvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter(makeConfig())

	req, _ := http.NewRequest("DELETE", "http://localhost/api/links/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func makeConfig() *Config {
	return NewConfig(false, "", "8080")
}
