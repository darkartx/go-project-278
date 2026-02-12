package code

import (
	"bytes"
	db "code/db/generated"
	"code/handlers"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

var conn *sql.DB

//go:embed db/migrations
var migrationsFS embed.FS

func TestPingRoute(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestLinksList(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		var err error
		var links [2]db.Link

		for i := 0; i < 2; i++ {
			links[i], err = q.CreateLink(ctx, db.CreateLinkParams{
				OriginalUrl: "https://google.com",
				ShortName:   fmt.Sprintf("test%d", i),
			})

			if err != nil {
				t.Fatalf("create link: %v", err)
			}
		}

		req, _ := http.NewRequest("GET", "http://localhost/api/links", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "links 0-9/2", w.Header().Get("Content-Range"))

		expectedLinks := []handlers.Link{
			{Id: uint64(links[0].ID), OriginalUrl: "https://google.com", ShortName: "test0", ShortUrl: "http://localhost/r/test0"},
			{Id: uint64(links[1].ID), OriginalUrl: "https://google.com", ShortName: "test1", ShortUrl: "http://localhost/r/test1"},
		}
		var actualLinks []handlers.Link
		err = json.Unmarshal(w.Body.Bytes(), &actualLinks)
		assert.NoError(t, err)
		assert.Equal(t, expectedLinks, actualLinks)
	})
}

func TestLinksListWithPagination(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		var err error
		var links [20]db.Link

		for i := 0; i < 20; i++ {
			links[i], err = q.CreateLink(ctx, db.CreateLinkParams{
				OriginalUrl: "https://google.com",
				ShortName:   fmt.Sprintf("test%d", i),
			})

			if err != nil {
				t.Fatalf("create link: %v", err)
			}
		}

		req, _ := http.NewRequest("GET", "http://localhost/api/links?range=[5,10]", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "links 5-10/20", w.Header().Get("Content-Range"))

		expectedLinks := make([]handlers.Link, 0, 6)

		for _, link := range links[5:11] {
			expectedLinks = append(
				expectedLinks,
				handlers.Link{
					Id:          uint64(link.ID),
					OriginalUrl: "https://google.com",
					ShortName:   link.ShortName,
					ShortUrl:    fmt.Sprintf("http://localhost/r/%s", link.ShortName),
				},
			)
		}

		var actualLinks []handlers.Link
		err = json.Unmarshal(w.Body.Bytes(), &actualLinks)
		assert.NoError(t, err)
		assert.Equal(t, expectedLinks, actualLinks)
	})
}

func TestLinksListWithInvalidPagination(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		var err error
		var links [2]db.Link

		for i := 0; i < 2; i++ {
			links[i], err = q.CreateLink(ctx, db.CreateLinkParams{
				OriginalUrl: "https://google.com",
				ShortName:   fmt.Sprintf("test%d", i),
			})

			if err != nil {
				t.Fatalf("create link: %v", err)
			}
		}

		cases := []string{
			"[abc,10]",
			"[-1,10]",
			"[20,10]",
			"asdasd",
		}

		for _, caseItem := range cases {
			req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost/api/links?range=%s", caseItem), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			expected := `{"error":"Bad Request","message":"invalid range param"}`
			assert.JSONEq(t, expected, w.Body.String())
		}
	})
}

func TestLinksCreate(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		body := `{"original_url":"https://google.com","short_name":"testtest"}`
		req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var actualLink handlers.Link
		var createdLink db.Link
		var err error

		err = json.Unmarshal(w.Body.Bytes(), &actualLink)
		assert.NoError(t, err)

		createdLink, err = q.GetLink(ctx, int64(actualLink.Id))

		if err != nil {
			t.Fatalf("get link: %v", err)
		}

		assert.Equal(t, actualLink.OriginalUrl, createdLink.OriginalUrl)
		assert.Equal(t, actualLink.OriginalUrl, "https://google.com")
		assert.Equal(t, actualLink.ShortName, createdLink.ShortName)
		assert.Equal(t, actualLink.ShortName, "testtest")
		assert.Equal(t, actualLink.ShortUrl, "http://localhost/r/testtest")
	})
}

func TestLinksCreateWithoutShortName(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		body := `{"original_url":"https://google.com"}`
		req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var actualLink handlers.Link
		var createdLink db.Link
		var err error

		err = json.Unmarshal(w.Body.Bytes(), &actualLink)
		assert.NoError(t, err)

		createdLink, err = q.GetLink(ctx, int64(actualLink.Id))

		if err != nil {
			t.Fatalf("get link: %v", err)
		}

		assert.Equal(t, "https://google.com", actualLink.OriginalUrl)
		assert.Equal(t, createdLink.ShortName, actualLink.ShortName)
		assert.Equal(t, "http://localhost/r/"+createdLink.ShortName, actualLink.ShortUrl)
	})
}

func TestLinksCreateWithInvalidOriginalUrl(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		body := `{"original_url":"invalid-url","short_name":"testtest"}`
		req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expected := `{"error":"Bad Request","message":"invalid original url"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestLinksCreateWithInvalidShortName(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		body := `{"original_url":"http://google.com","short_name":"!@#$!asdasd"}`
		req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expected := `{"error":"Bad Request","message":"invalid short name"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestLinksCreateWithUsedShortName(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		if _, err := q.CreateLink(ctx, db.CreateLinkParams{OriginalUrl: "https://google.com", ShortName: "testtest"}); err != nil {
			t.Fatalf("create link: %v", err)
		}

		body := `{"original_url":"https://google.com","short_name":"testtest"}`
		req, _ := http.NewRequest("POST", "http://localhost/api/links", bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		expected := `{"error":"Unprocessable Entity","message":"short name is already used"}`
		assert.JSONEq(t, expected, w.Body.String())
	})
}

func TestLinksGet(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		link, err := q.CreateLink(ctx, db.CreateLinkParams{OriginalUrl: "https://google.com", ShortName: "testtest"})
		if err != nil {
			t.Fatalf("create link: %v", err)
		}

		req, _ := http.NewRequest("GET", fmt.Sprint("http://localhost/api/links/", link.ID), nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var actualLink handlers.Link
		err = json.Unmarshal(w.Body.Bytes(), &actualLink)
		assert.NoError(t, err)

		assert.Equal(t, link.ID, int64(actualLink.Id))
		assert.Equal(t, "https://google.com", actualLink.OriginalUrl)
		assert.Equal(t, "testtest", actualLink.ShortName)
		assert.Equal(t, "http://localhost/r/testtest", actualLink.ShortUrl)
	})
}

func TestLinksGetWithNotExistingId(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "http://localhost/api/links/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := `{"error":"Not Found","message":"Not found"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksGetWithInvalidId(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "http://localhost/api/links/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdate(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		link, err := q.CreateLink(ctx, db.CreateLinkParams{OriginalUrl: "http://localhost/", ShortName: "123ABC"})
		if err != nil {
			t.Fatalf("create link: %v", err)
		}

		body := `{"original_url":"https://google.com","short_name":"testtest"}`
		req, _ := http.NewRequest("PUT", fmt.Sprint("http://localhost/api/links/", link.ID), bytes.NewBufferString(body))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		link, err = q.GetLink(ctx, link.ID)
		if err != nil {
			t.Fatalf("get link: %v", err)
		}

		var actualLink handlers.Link
		err = json.Unmarshal(w.Body.Bytes(), &actualLink)
		assert.NoError(t, err)

		assert.Equal(t, link.ID, int64(actualLink.Id))
		assert.Equal(t, "https://google.com", actualLink.OriginalUrl)
		assert.Equal(t, "https://google.com", link.OriginalUrl)
		assert.Equal(t, "testtest", actualLink.ShortName)
		assert.Equal(t, "testtest", link.ShortName)
		assert.Equal(t, "http://localhost/r/testtest", actualLink.ShortUrl)
	})
}

func TestLinksUpdateWithInvalidId(t *testing.T) {
	router := setupTestRouter()

	body := `{"original_url":"http://google.com","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/abc", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithInvalidOriginalUrl(t *testing.T) {
	router := setupTestRouter()

	body := `{"original_url":"invalid-url","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	expected := `{"error":"Unprocessable Entity","message":"invalid original url"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithInvalidShortName(t *testing.T) {
	router := setupTestRouter()

	body := `{"original_url":"http://google.com","short_name":"!@#$!asdasd"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	expected := `{"error":"Unprocessable Entity","message":"invalid short name"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksUpdateWithNotExistingId(t *testing.T) {
	router := setupTestRouter()

	body := `{"original_url":"https://google.com","short_name":"testtest"}`
	req, _ := http.NewRequest("PUT", "http://localhost/api/links/1", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := `{"error":"Not Found","message":"Not found"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksDelete(t *testing.T) {
	withTx(t, func(ctx context.Context, q *db.Queries, tx *sql.Tx) {
		router := setupTestRouterWithQueries(q)

		link, err := q.CreateLink(ctx, db.CreateLinkParams{OriginalUrl: "http://localhost/", ShortName: "123ABC"})
		if err != nil {
			t.Fatalf("create link: %v", err)
		}

		req, _ := http.NewRequest("DELETE", fmt.Sprint("http://localhost/api/links/", link.ID), nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "", w.Body.String())

		_, err = q.GetLink(ctx, link.ID)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestLinksDeleteWithInvalidId(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "http://localhost/api/links/abc", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Bad Request","message":"invalid id"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestLinksDeleteWithNotExistingId(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "http://localhost/api/links/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := `{"error":"Not Found","message":"Not found"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	databaseUrl := os.Getenv("DATABASE_URL")

	conn, err = sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	if err := conn.PingContext(ctxPing); err != nil {
		cancel()
		log.Fatalf("ping db: %v", err)
	}
	cancel()

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect: %v", err)
	}

	if err := goose.Up(conn, "db/migrations"); err != nil {
		log.Fatalf("goose up: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

func withTx(t *testing.T, fn func(ctx context.Context, q *db.Queries, tx *sql.Tx)) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	t.Cleanup(cancel)

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}

	t.Cleanup(func() { _ = tx.Rollback() })

	qtx := db.New(tx)
	fn(ctx, qtx, tx)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return setupRouter(db.New(conn), NewConfig(false, "", "8080"))
}

func setupTestRouterWithQueries(queries *db.Queries) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return setupRouter(queries, NewConfig(false, "", "8080"))
}
