package integration

import (
	"avito2024/internal/app"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetHouseFlatsService(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16.4-alpine3.20",
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "migrations", "000001_init.up.sql")),
		postgres.WithInitScripts(filepath.Join("addData.sql")),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	os.Setenv("DATABASE_URL", connStr)

	gin.SetMode(gin.TestMode)
	r := app.SetupRouter()

	go func() {
		r.Run(":8080")
	}()

	tests := []struct {
		name               string
		id                 int
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			name:               "Unauthorized request (missing cookies)",
			id:                 1,
			expectedHTTPStatus: 401,
			expectedResponse:   "",
		},
		{
			name:               "Valid response for existing house",
			id:                 1,
			expectedHTTPStatus: 200,
			expectedResponse: `[
			{
				"id": 2,
				"house_id": 1,
				"price": 250000,
				"rooms": 2,
				"status": "approved"
			},
			{
				"id": 5,
				"house_id": 1,
				"price": 500000,
				"rooms": 5,
				"status": "approved"
			}
			]`,
		},
		{
			name:               "House does not exist",
			id:                 999,
			expectedHTTPStatus: 200,
			expectedResponse:   "[]",
		},
		{
			name:               "Moderator - mixed status flats",
			id:                 1,
			expectedHTTPStatus: 200,
			expectedResponse: `[
			{
				"id": 1,
				"house_id": 1,
				"price": 300000,
				"rooms": 3,
				"status": "created"
			},
			{
				"id": 2,
				"house_id": 1,
				"price": 250000,
				"rooms": 2,
				"status": "approved"
			},
			{
				"id": 3,
				"house_id": 1,
				"price": 400000,
				"rooms": 4,
				"status": "on moderation"
			},
			{
				"id": 4,
				"house_id": 1,
				"price": 150000,
				"rooms": 1,
				"status": "declined"
			},
			{
				"id": 5,
				"house_id": 1,
				"price": 500000,
				"rooms": 5,
				"status": "approved"
			},
			{
				"id": 6,
				"house_id": 1,
				"price": 350000,
				"rooms": 3,
				"status": "created"
			}
			]`,
		},
		{
			name:               "Empty flats list",
			id:                 102,
			expectedHTTPStatus: 200,
			expectedResponse:   "[]",
		},
		{
			name:               "Invalid house ID",
			id:                 -1,
			expectedHTTPStatus: 400,
			expectedResponse:   "",
		},
	}

	t.Run(tests[0].name, func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/house/%d", tests[0].id), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, tests[0].expectedHTTPStatus, w.Code)

		var actualResponse string
		bodyBytes, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		actualResponse = string(bodyBytes)

		assert.Equal(t, tests[0].expectedResponse, actualResponse)
	})

	loginReq, _ := http.NewRequest("GET", "/dummyLogin?user_type=client", nil)
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	client := loginW.Result().Cookies()

	loginReq, _ = http.NewRequest("GET", "/dummyLogin?user_type=moderator", nil)
	loginW = httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)

	moderator := loginW.Result().Cookies()

	for i := 1; i < len(tests); i++ {
		i := i
		t.Run(tests[i].name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/house/%d", tests[i].id), nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			if i != 3 {
				req.AddCookie(client[0])
			} else {
				req.AddCookie(moderator[0])
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tests[i].expectedHTTPStatus, w.Code)

			var actualResponse string
			bodyBytes, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}
			actualResponse = string(bodyBytes)

			if json.Valid([]byte(tests[i].expectedResponse)) {
				assert.JSONEq(t, tests[i].expectedResponse, actualResponse)
			} else {
				assert.Equal(t, tests[i].expectedResponse, actualResponse)
			}
		})
	}

}
