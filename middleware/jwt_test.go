package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyJWT(t *testing.T) {
	token := GenerateJWT("my_secret_key", "test")
	verify := VerifyJWT("my_secret_key")
	r := gin.Default()
	r.GET("/auth", verify)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	r.ServeHTTP(w, req)
	if w.Code != 403 {
		t.Errorf("Return %v without auth", w.Code)
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/auth", nil)
	req.Header.Set("authorization", "Bearer " + token)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Return %v with auth %v", w.Code, w.Body)
	}
}
