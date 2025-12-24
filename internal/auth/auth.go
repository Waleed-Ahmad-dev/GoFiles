package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"GoFiles/internal/config"
	"GoFiles/internal/types"
	"GoFiles/internal/utils"

	"github.com/google/uuid"
)

var sessions = map[string]string{}

// HandleSystemStatus tells the Frontend if we need Setup or Login
func HandleSystemStatus(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	if config.IsConfigured {
		w.Write([]byte(`{"status": "ready"}`)) // Show Login Screen
	} else {
		w.Write([]byte(`{"status": "setup_required"}`)) // Show Setup Screen
	}
}

// HandleSetup is the "First Run" wizard
func HandleSetup(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	if config.IsConfigured {
		http.Error(w, "System is already configured", http.StatusForbidden)
		return
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and Password required", http.StatusBadRequest)
		return
	}

	// Save to gofiles.json
	if err := config.SaveConfig(req.Username, req.Password); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	config.IsConfigured = true // Switch to Normal Mode

	// Auto-login the user
	createSession(w, req.Username)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Setup complete"}`))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method != http.MethodPost {
		return
	}

	// If setup isn't done, we can't login!
	if !config.IsConfigured {
		http.Error(w, "Setup required first", http.StatusLocked)
		return
	}

	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Check credentials against Loaded Config
	if req.Username == config.AppConfig.Username && req.Password == config.AppConfig.Password {
		createSession(w, req.Username)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login successful"}`))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	c, err := r.Cookie("session_token")
	if err == nil {
		delete(sessions, c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: "", Expires: time.Now().Add(-1 * time.Hour), HttpOnly: true, Path: "/",
	})
	w.WriteHeader(http.StatusOK)
}

func HandleCheckAuth(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
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
		utils.EnableCors(&w)
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		// 1. If Setup is NOT done, block everything except setup/status endpoints
		if !config.IsConfigured {
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
