package ratelimitermiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/rcdmk/go-ratelimiter/ratelimitermiddleware"
)

func Test_StdLib_Allows_Requests_For_Different_Sources(t *testing.T) {
	tests := []struct {
		name           string
		headerValue    string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			headerValue:    "test1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request 2",
			headerValue:    "test2",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request 3",
			headerValue:    "test3",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request 4",
			headerValue:    "test4",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request 5",
			headerValue:    "test5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request 6",
			headerValue:    "test6",
			expectedStatus: http.StatusOK,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	options := Options{
		MaxRatePerSecond: 10,
		MaxBurst:         5,
		SourceHeaderKey:  "Authorization",
	}

	middleware := StdLib(handler, options)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(options.SourceHeaderKey, tt.headerValue)

			res := httptest.NewRecorder()

			middleware.ServeHTTP(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, res.Code)
			}
		})
	}
}

func Test_StdLib_Allows_Burst_Requests_For_Multiple_Sources(t *testing.T) {
	tests := []struct {
		name                   string
		requestCount           int
		expectedLastStatus     int
		expectedRequestsServed int
	}{
		{
			name:                   "Valid burst 1",
			requestCount:           5,
			expectedLastStatus:     http.StatusOK,
			expectedRequestsServed: 5,
		},
		{
			name:                   "Valid burst 2",
			requestCount:           10,
			expectedLastStatus:     http.StatusOK,
			expectedRequestsServed: 10,
		},
		{
			name:                   "Invalid burst 1",
			requestCount:           11,
			expectedLastStatus:     http.StatusTooManyRequests,
			expectedRequestsServed: 10, // last will be rejected
		},
	}

	handlerCalled := 0
	headerKey := "Authorization"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		handlerCalled++
	})

	options := Options{
		MaxRatePerSecond: 5,
		MaxBurst:         10,
		SourceHeaderKey:  headerKey,
	}

	middleware := StdLib(handler, options)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set(headerKey, "test"+tt.name)

				res := httptest.NewRecorder()

				middleware.ServeHTTP(res, req)

				if i == tt.requestCount-1 {
					if res.Code != tt.expectedLastStatus {
						t.Errorf("Expected status code %d, but got %d", tt.expectedLastStatus, res.Code)
					}
				} else {
					if res.Code != http.StatusOK {
						t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
					}
				}
			}

			if handlerCalled != tt.expectedRequestsServed {
				t.Errorf("Expected handler to be called %d times, but got %d", tt.expectedRequestsServed, handlerCalled)
			}

			handlerCalled = 0
		})
	}

}
