package config

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

var (
	Store = sessions.NewCookieStore([]byte("your-secret-key"))
)

func InitSession() {
	gob.Register(time.Time{})
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}
