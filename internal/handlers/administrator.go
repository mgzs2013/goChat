package handlers

import (
	"fmt"
	"net/http"
)

// AdminHandler handles requests to the /admin endpoint
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome, Admin! You have access to the admin area.")
}
