package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid" // You will need to install this: go get github.com/google/uuid
)

// Store active sessions in memory (Simple and fast)
// Map: SessionToken -> Username
var sessions = map[string]string{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 1. Handle Login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Check credentials against config
	if req.Username == AdminUser && req.Password == AdminPass {
		// Generate a random session token
		token := uuid.New().String()

		// Save session
		sessions[token] = req.Username

		// Set Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true, // Prevents JavaScript from reading it (XSS protection)
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login successful"}`))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

// 2. Handle Logout
func handleLogout(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	c, err := r.Cookie("session_token")
	if err == nil {
		delete(sessions, c.Value) // Remove from server memory
	}

	// Expire the cookie immediately
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}

// 3. Middleware (The Security Guard)
// This wraps your other functions. It checks the cookie BEFORE letting them pass.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Allow OPTIONS requests (needed for CORS pre-flight checks)
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		// Get the cookie
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Check if session exists in our memory map
		sessionToken := c.Value
		user, exists := sessions[sessionToken]
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Success! Allow request to proceed
		// (Optional: You could pass 'user' context here if needed)
		fmt.Printf("User %s accessed %s\n", user, r.URL.Path)
		next(w, r)
	}
}

// CheckAuthStatus handles the /api/me endpoint
func handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// If the request reached here, it passed the Middleware, so they are logged in.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"authenticated": true}`))
}
