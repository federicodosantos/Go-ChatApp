package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/federicodosantos/Go-ChatApp/internal/user/usecase"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	usecase usecase.UserUCItf
	store *sessions.CookieStore
	logger *zap.Logger
}

func NewUserHandler(usecase usecase.UserUCItf, 
	store *sessions.CookieStore,
	logger *zap.Logger) *UserHandler {
		return &UserHandler{
			usecase: usecase,
			store: store,
			logger: logger,
		}
	}

func (uh *UserHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := uh.store.Get(r, "oauth-session")
	if err != nil {
		uh.logger.Error("Failed to get session", zap.Error(err))
		return
	}

	state := uuid.New().String()
	verifier := oauth2.GenerateVerifier()

	session.Values["state"] = state
	session.Values["verifier"] = verifier
	session.Save(r, w)

	url, err := uh.usecase.GoogleLogin(state, verifier)
	if err != nil {
		uh.logger.Error("cannot login with google", zap.Error(err))
		return 
	}

	uh.logger.Info("Redirect to Auth URL", zap.String("Auth Url", url))

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (uh *UserHandler) CallBackGoogle(w http.ResponseWriter, r *http.Request) {
	uh.logger.Info("Callback")

	session, err := uh.store.Get(r, "oauth-session")
	if err != nil {
		uh.logger.Error("Failed to get session", zap.Error(err))
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	state := r.FormValue("state")
	sessionState, ok := session.Values["state"].(string)
	if !ok || state != sessionState {
		uh.logger.Info("Invalid OAuth state", zap.String("received_state", state), zap.String("expected_state", sessionState))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	authCode := r.FormValue("code")
	if authCode == "" {
		uh.logger.Warn("Auth Code not found")
		w.Write([]byte("Auth Code not found to provide AccessToken"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		return
	}
	
	verifier, ok := session.Values["verifier"].(string)
	uh.logger.Info("verifier code in handler", zap.String("verifier code" ,verifier))
	if !ok {
		uh.logger.Error("Verifier not found in session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := uh.usecase.ExchangeToken(authCode, verifier)
	if err != nil {
		if oauthErr, ok := err.(*oauth2.RetrieveError); ok {
			uh.logger.Error("OAuth error details",
				zap.Int("StatusCode", oauthErr.Response.StatusCode),
				zap.String("Body", string(oauthErr.Body)))
			http.Error(w, "Failed to exchange token", oauthErr.Response.StatusCode)
		} else {
			uh.logger.Error("Failed to exchange token", zap.Error(err))
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		}
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		uh.logger.Error("Error to fetch user data", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		uh.logger.Error("Failed to parse JSON", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token.SetAuthHeader(r)

	var googleUserInfo map[string]interface{}
	err = json.Unmarshal(userData, &googleUserInfo)
	if err != nil {
		uh.logger.Error("Failed to unmarshal user data", zap.Error(err))
		return
	}

	createdUser, err := uh.usecase.FetchUserData(googleUserInfo)
	if err != nil {
		uh.logger.Error("Cannot create user", zap.Error(err), zap.Any("user data" ,googleUserInfo))
		http.Error(w, "failed to parse user data", http.StatusInternalServerError)
		return
	}

	user, err := json.Marshal(createdUser)
	if err != nil {
		uh.logger.Error("Cannot marshal user data", zap.Error(err))
		return
	}

	w.Write([]byte("Success boskuh...\n"))
	w.Write([]byte(string(user)))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}