package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var sessions = map[string]string{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// handleSystemStatus tells the Frontend if we need Setup or Login
func handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	if IsConfigured {
		w.Write([]byte(`{"status": "ready"}`)) // Show Login Screen
	} else {
		w.Write([]byte(`{"status": "setup_required"}`)) // Show Setup Screen
	}
}

// handleSetup is the "First Run" wizard
func handleSetup(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	if IsConfigured {
		http.Error(w, "System is already configured", http.StatusForbidden)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password required", http.StatusBadRequest)
		return
	}

	// Save to gofiles.json
	if err := SaveConfig(req.Username, req.Password); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	IsConfigured = true // Switch to Normal Mode

	// Auto-login the user
	createSession(w, req.Username)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Setup complete"}`))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	// If setup isn't done, we can't login!
	if !IsConfigured {
		http.Error(w, "Setup required first", http.StatusLocked)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Check credentials against Loaded Config
	if req.Username == AppConfig.Username && req.Password == AppConfig.Password {
		createSession(w, req.Username)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login successful"}`))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	c, err := r.Cookie("session_token")
	if err == nil {
		delete(sessions, c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: "", Expires: time.Now().Add(-1 * time.Hour), HttpOnly: true, Path: "/",
	})
	w.WriteHeader(http.StatusOK)
}

func handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"authenticated": true}`))
}

// --- HELPER ---
func createSession(w http.ResponseWriter, username string) {
	token := uuid.New().String()
	sessions[token] = username
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: token, Expires: time.Now().Add(24 * time.Hour), HttpOnly: true, Path: "/",
	})
}

// --- MIDDLEWARE ---
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		// 1. If Setup is NOT done, block everything except setup/status endpoints
		if !IsConfigured {
			http.Error(w, "Setup Required", http.StatusLocked) // 423 Locked
			return
		}

		// 2. Normal Auth Check
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		sessionToken := c.Value
		user, exists := sessions[sessionToken]
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Printf("User %s accessed %s\n", user, r.URL.Path)
		next(w, r)
	}
}