package handlers

import (
	"net/http"

	"github.com/fujidaiti/bookings/pkg/renderer"
)

func Home(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "home")
}

func About(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "about")
}
