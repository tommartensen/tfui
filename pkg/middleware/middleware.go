package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/tommartensen/tfui/pkg/config"
)

func authenticateRequest(r *http.Request) error {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return fmt.Errorf("no authorization header passed")
	}
	splitToken := strings.Split(auth, " ")
	if len(splitToken) != 2 {
		return fmt.Errorf("authorization header malformed '%s'", auth)
	}
	authType := strings.TrimSpace(splitToken[0])
	if authType != "Bearer" {
		return fmt.Errorf("only authorization header type 'Bearer' is supported, not '%s'", authType)
	}
	b64Token := strings.TrimSpace(splitToken[1])
	token, err := base64.StdEncoding.DecodeString(b64Token)
	if err != nil {
		return fmt.Errorf("authorization header token base64 malformed '%s'", b64Token)
	}

	cfg := config.New()
	if string(token) != cfg.ApplicationToken {
		return errors.New("secret mismatch")
	}
	return nil
}

func logRequest(r *http.Request) {
	log.Printf("Received call [%s] %s\n", r.Method, strings.Replace(r.URL.Path, "\n", "", -1))
}

// Middleware function called for each request
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		if err := authenticateRequest(r); err != nil {
			log.Printf("Authentication failed: %s", err.Error())
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
