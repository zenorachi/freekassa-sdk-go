package freekassa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCallbacks_Middleware(t *testing.T) {
	tests := []struct {
		name           string
		whitelist      map[string]struct{}
		remoteAddr     string
		expectedStatus int
	}{
		{
			name: "AllowsWhitelistedIP",
			whitelist: map[string]struct{}{
				"1.2.3.4": {},
			},
			remoteAddr:     "1.2.3.4:12345",
			expectedStatus: http.StatusOK,
		},
		{
			name: "BlocksNonWhitelistedIP",
			whitelist: map[string]struct{}{
				"1.2.3.4": {},
			},
			remoteAddr:     "5.6.7.8:12345",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			c.Request.RemoteAddr = tt.remoteAddr

			router := gin.New()
			router.Use(middleware(tt.whitelist))
			router.GET("/", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			router.ServeHTTP(w, c.Request)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
