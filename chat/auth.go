package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
)

// ChatUser represents a user of the chat service.
type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated, redirect
		// TODO: Shit authentication mechanism.
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	// on success call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth wraps a HttpHandler with authentication support; the wrapped handler will not be called if authentication
// fails.
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// TODO: use a library to do URL matching instead of trying to split URLs and handle shit in switch / cases
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		fmt.Println("Beggining auth for", provider)
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
			return
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error trying to GetBeginAuthURL for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
			return
		}

		m := md5.New()
		_, err = io.WriteString(m, strings.ToLower(user.Email()))

		if err != nil {
			// TODO: Should do a non user related error handler.
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		userID := fmt.Sprintf("%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     userID,
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
			"email":      user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		// TODO: XSS vulnerability
		fmt.Fprintf(w, "Auth action %s is not supported", action)
	}
}
