package integration

import (
	"avito2024/internal/app"
	"bytes"
	"context"
	"encoding/json"
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

func TestCreateFlatService(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:14.13-alpine3.20",
		postgres.WithInitScripts(filepath.Join("..", "..", "..", "migrations", "init.sql")),
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
		houseId            int
		price              int
		rooms              int
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			name:               "Unauthorized request (missing cookies)",
			houseId:            1,
			price:              1,
			rooms:              1,
			expectedHTTPStatus: 401,
			expectedResponse:   "",
		},
		{
			name:               "Valid request",
			houseId:            1,
			price:              100,
			rooms:              10,
			expectedHTTPStatus: 200,
			expectedResponse:   "{\n  \"house_id\": 1,\n  \"id\": 7,\n  \"price\": 100,\n  \"rooms\": 10,\n  \"status\": \"created\"\n}",
		},
		{
			name:               "Invalid house ID",
			houseId:            102,
			price:              0,
			rooms:              0,
			expectedHTTPStatus: 400,
			expectedResponse:   "",
		},
		{
			name:               "Invalid price",
			houseId:            1,
			price:              -1,
			rooms:              1,
			expectedHTTPStatus: 400,
			expectedResponse:   "",
		},
		{
			name:               "Invalid rooms",
			houseId:            1,
			price:              1,
			rooms:              -1,
			expectedHTTPStatus: 400,
			expectedResponse:   "",
		},
	}

	t.Run(tests[0].name, func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"house_id": tests[0].houseId,
			"price":    tests[0].price,
			"rooms":    tests[0].rooms,
		})
		req, _ := http.NewRequest("POST", "/flat/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
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
			reqBody, _ := json.Marshal(map[string]interface{}{
				"house_id": tests[i].houseId,
				"price":    tests[i].price,
				"rooms":    tests[i].rooms,
			})
			req, err := http.NewRequest("POST", "/flat/create", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

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
