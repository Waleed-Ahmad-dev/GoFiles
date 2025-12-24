package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitTrash()

	// --- PUBLIC ROUTES ---
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)

	// --- PROTECTED ROUTES (Wrapped in AuthMiddleware) ---

	// Utility to check if user is logged in
	http.HandleFunc("/api/me", AuthMiddleware(handleCheckAuth))

	// Read & Search
	http.HandleFunc("/api/files", AuthMiddleware(handleListFiles))
	http.HandleFunc("/api/download", AuthMiddleware(handleDownloadFile))
	http.HandleFunc("/api/search", AuthMiddleware(handleSearch))

	// Trash
	http.HandleFunc("/api/trash/list", AuthMiddleware(handleListTrash))
	http.HandleFunc("/api/trash/restore", AuthMiddleware(handleRestore))
	http.HandleFunc("/api/trash/empty", AuthMiddleware(handleEmptyTrash))

	// Write
	http.HandleFunc("/api/upload", AuthMiddleware(handleUploadFile))
	http.HandleFunc("/api/mkdir", AuthMiddleware(handleCreateDir))
	http.HandleFunc("/api/delete", AuthMiddleware(handleDelete))

	// Organize
	http.HandleFunc("/api/rename", AuthMiddleware(handleRename))
	http.HandleFunc("/api/move", AuthMiddleware(handleMove))
	http.HandleFunc("/api/copy", AuthMiddleware(handleCopy))

	fmt.Println("🚀 GoFiles Server started on http://localhost:8080")
	fmt.Printf("🔑 Default Login: %s / %s\n", AdminUser, AdminPass)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
