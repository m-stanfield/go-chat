package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (s *Server) logEndpoint(next http.Handler) http.Handler {
	counter := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter = counter + 1
		start_time := time.Since(startTime)
		// Proceed with the next handler
		next.ServeHTTP(w, r)
		end_time := time.Since(startTime.Add(start_time))

		log.Printf(
			"%d Endpoint hit: %s took %d ms\n",
			counter,
			r.URL,
			end_time.Milliseconds(),
		)
	})
}

func (s *Server) WithAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieName := "token"
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				// Handle the case where the cookie is not found
				http.Error(w, "Token cookie not found", http.StatusUnauthorized)
				return
			}
			// Handle other potential errors
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			return
		}

		// Access the cookie value
		token := cookie.Value
		passwordInfo, err := s.db.GetUserLoginInfoFromToken(token)
		if err != nil {
			http.Error(w, "unable to locate password", http.StatusBadRequest)
			return
		}

		if !s.validSession(passwordInfo, token) {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), "userid", passwordInfo.UserId)))
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().
			Set("Access-Control-Allow-Origin", "*")
			// Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().
			Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().
			Set("Access-Control-Allow-Credentials", "true")
			// Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
