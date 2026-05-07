package session

import (
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func Init() {
	if sessionManager != nil {
		panic("Init must be called only once.")
	}

	sessionManager = scs.New()
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Lifetime = 48 * time.Hour
}

// LoadAndSave returns a middleware to load and save the session data for the given request.
func LoadAndSave() func(http.Handler) http.Handler {
	return sessionManager.LoadAndSave
}

const guestIdKey = "guestId"

func SetGuestCredential(id int, r *http.Request) {
	ctx := r.Context()
	sessionManager.Destroy(ctx)
	sessionManager.Put(ctx, guestIdKey, id)
}

func GetGuestCredential(r *http.Request) (int, error) {
	id := sessionManager.GetInt(r.Context(), guestIdKey)
	if id == 0 {
		return 0, errors.New("Guest id is not found in the session data.")
	}
	return id, nil
}

func IsGuestLoggedIn(r *http.Request) bool {
	_, err := GetGuestCredential(r)
	return err == nil
}
