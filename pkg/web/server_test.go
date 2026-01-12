//go:build test

package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/platform/ctxval"

	"adeynack.net/lapiasse/pkg/web"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Those are basic SMOKE TESTS to ensure that the web server is correctly
// wired to the API implementation and that errors are correctly handled
// end-to-end. Detailed tests of the API implementation are done in the
// controller package.
func TestServer(t *testing.T) {
	for name, tc := range map[string]struct {
		init           func(ctx context.Context) context.Context
		urlPath        string
		method         string
		requestBody    string
		expectedStatus int
		expectedBody   string
		auth           func() (tokenValue string, includeAuthorizationHeader bool) // by default, valid token is provided
	}{
		"not providing an Authorization header": {
			urlPath:        "/health",
			method:         http.MethodGet,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"type":"error:unauthorized","title":"Unauthorized","status":401}`,
			auth:           func() (string, bool) { return "", false },
		},
		"providing an invalid API token": {
			urlPath:        "/health",
			method:         http.MethodGet,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"type":"error:unauthorized","title":"Unauthorized","status":401}`,
			auth:           func() (string, bool) { return "invalid-token", true },
		},
		"a path returning a simple response - GET health": {
			urlPath:        "/health",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"healthy"}`,
		},
		"a non existing path - GET nonexistent-path": {
			urlPath:        "/nonexistent-path",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"type":"error:not_found","title":"Not Found","status":404,"detail":"Path \"/nonexistent-path\" not found"}`,
		},
		"an existing path on a non supported method - PATCH health": {
			urlPath:        "/health",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   `{"type":"error:method_not_allowed","title":"Method Not Allowed","status":405,"detail":"Path \"/health\" does not allow verb \"PATCH\""}`,
		},
		"a path querying the database - GET books": {
			urlPath:        "/books",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"books":[],"pagination":{"has_next":false}}`,
		},
		"a path returning a Not Found error - GET books nonexistent-id": {
			urlPath:        "/books/nonexistent-id",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"type":"error:not_found","title":"Not Found","status":404,"detail":"Book with ID \"nonexistent-id\" not found"}`,
		},
		"a POST without a JSON request body - POST books": {
			urlPath:        "/books",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"type":"error:bad_request","title":"Bad Request","status":400,"detail":"can't decode JSON body: EOF"}`,
		},
		"a path returning an Unprocessable Entity error - POST books with invalid data": {
			urlPath:        "/books",
			method:         http.MethodPost,
			requestBody:    `{"book": {"name": ""}}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: `{
				"status": 422,
				"title": "Resource did not validate",
				"type": "error:validation",
				"validation_errors": [
					{
						"field": "Book.Name",
						"message": "is required",
						"validation": "required"
					},
					{
						"field": "Book.DefaultCurrencyIsoCode",
						"message": "is required",
						"validation": "required"
					}
				]
			}`,
		},
		"panic during request handling returns JSON error": {
			init: func(ctx context.Context) context.Context {
				return ctxval.Register[*gorm.DB](ctx, nil) // nil DB will cause a panic
			},
			urlPath:        "/books",
			method:         http.MethodGet,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
				"status": 500,
				"title": "Internal Server Error",
				"type": "error:internal_error"
			}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)
			if tc.init != nil {
				ctx = tc.init(ctx)
			}

			config, err := web.ConfigurationDefaults()
			require.NoError(t, err)

			config.Expose = true
			config.Port = 0 // use a random available port

			server, err := web.StartServer(ctx, config)
			require.NoError(t, err)
			t.Cleanup(func() {
				lo.Must0(server.Close())
			})

			request := httptest.NewRequest(tc.method, tc.urlPath, strings.NewReader(tc.requestBody))

			if tc.auth == nil {
				request.Header.Set("Authorization", "Bearer "+web.StaticApiTokenTest)
			} else if expectedToken, includeAuth := tc.auth(); includeAuth {
				request.Header.Set("Authorization", "Bearer "+expectedToken)
			}

			recorder := httptest.NewRecorder()

			server.Handler.ServeHTTP(recorder, request)

			responseBody := recorder.Body.String()
			response := recorder.Result()

			require.Equalf(t, tc.expectedStatus, response.StatusCode, "Body:\n%s", responseBody)
			require.Equal(t, "application/json", response.Header.Get("Content-Type"), "Response:\n%s", responseBody)

			require.JSONEqf(t, tc.expectedBody, responseBody, "Body:\n%s", responseBody)
		})
	}
}
