package handlers

import (
	"fmt"
	"net/http"
)

// EditorHandler handles requests to the /editor endpoint
func EditorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome, Editor! You have access to the editor area.")
}
