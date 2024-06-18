package google

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/gorilla/sessions"
)

type OauthGoogle struct {
	conf *oauth2.Config
	logger *zap.Logger
	store *sessions.CookieStore
}

func NewOauthGoogle(conf *oauth2.Config, logger *zap.Logger, store *sessions.CookieStore) *OauthGoogle {
	return &OauthGoogle{
		conf: conf,
		logger: logger,
		store: store,
	}
}

// initialize oauth goole function
func InitOauthGoogle() *oauth2.Config {
	return &oauth2.Config{
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL: os.Getenv("GOOGLE_CALLBACK"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
}

func (oauth *OauthGoogle) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := oauth.store.Get(r, "oauth-session")
	if err != nil {
		oauth.logger.Error("Failed to get session", zap.Error(err))
		return
	}

	state := uuid.New().String()
	verifier := oauth2.GenerateVerifier()

	session.Values["state"] = state
	session.Values["verifier"] = verifier
	session.Save(r, w)

	url := oauth.conf.AuthCodeURL(state, 
	oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	oauth.logger.Info("Redirect to Auth URL", zap.String("Auth Url", url))

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (oauth *OauthGoogle) CallBackGoogle(w http.ResponseWriter, r *http.Request) {
	oauth.logger.Info("Callback")

	session, err := oauth.store.Get(r, "oauth-session")
	if err != nil {
		oauth.logger.Error("Failed to get session", zap.Error(err))
		return 
	}

	state := r.FormValue("state")
	sessionState := session.Values["state"].(string)
	if state != sessionState {
		oauth.logger.Info("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	if code == "" {
		oauth.logger.Warn("Code not found")
		w.Write([]byte("Code not found to provide AccessToken"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	} else {
		verifier := session.Values["verifier"].(string)
		token, err := oauth.conf.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			oauth.logger.Error("Exchange failed", zap.Error(err))
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			oauth.logger.Error("Error to fetch data user", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return 
		}
		defer resp.Body.Close()

		userData, err := io.ReadAll(resp.Body)
		if err != nil {
			oauth.logger.Error("Failed to parse JSON", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return 
		}

		w.Write([]byte("Success boskuh...\n"))
		w.Write([]byte(string(userData)))
		return
	}
}

