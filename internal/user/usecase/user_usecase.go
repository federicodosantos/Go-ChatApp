package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/federicodosantos/Go-ChatApp/internal/user"
	"github.com/federicodosantos/Go-ChatApp/internal/user/repository"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type UserUCItf interface {
	GoogleLogin(state, verifierCode string) (string, error)
	FetchUserData(userData map[string]interface{}) (*user.User, error)
	ExchangeToken(authCode, verifierCode string) (*oauth2.Token, error)
}

type UserUC struct {
	userRepo repository.UserRepoItf
	oauth    *oauth2.Config
	logger   *zap.Logger
}


func NewUserUC(userRepo repository.UserRepoItf,
	oauth *oauth2.Config,
	logger *zap.Logger) UserUCItf {
	return UserUC{userRepo: userRepo, oauth: oauth, logger: logger}
}

// Login implements UserUCItf.
func (u UserUC) FetchUserData(userData map[string]interface{}) (*user.User, error) {
	if userData == nil {
		return nil, errors.New("user data is nil")
	}

	var id, name, email, picture string
    var idOk, nameOk, emailOk, pictureOk bool

    if id, idOk = userData["id"].(string); !idOk {
        u.logger.Error("Failed to assert id to string")
        return nil, errors.New("id is not a string")
    }
    if name, nameOk = userData["name"].(string); !nameOk {
        u.logger.Error("Failed to assert name to string")
        return nil, errors.New("name is not a string")
    }
    if email, emailOk = userData["email"].(string); !emailOk {
        u.logger.Error("Failed to assert email to string")
        return nil, errors.New("email is not a string")
    }
    if picture, pictureOk = userData["picture"].(string); !pictureOk {
        u.logger.Error("Failed to assert picture to string")
        return nil, errors.New("picture is not a string")
    }

    // Membuat user baru
    newUser := &user.User{
        ID:    id,
        Name:  name,
        Email: email,
        Photo_Link: sql.NullString{
            String: picture,
            Valid:  true,
        },
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

	err := u.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// ExchangeToken implements UserUCItf.
func (u UserUC) ExchangeToken(authCode, verifierCode string) (*oauth2.Token, error) {
	if verifierCode == "" {
		return nil, errors.New("verifier code is nil")
	}

	token, err := u.oauth.Exchange(context.Background(),
		authCode, oauth2.VerifierOption(verifierCode))
	if err != nil {
		if oauthErr, ok := err.(*oauth2.RetrieveError); ok {
            u.logger.Error("OAuth error details",
                zap.Int("StatusCode", oauthErr.Response.StatusCode),
                zap.String("Body", string(oauthErr.Body)))
        }
        return nil, err
	}

	return token, nil
}

// GoogleLogin implements UserUCItf.
func (u UserUC) GoogleLogin(state string, verifierCode string) (string, error) {
	if  state == ""  {
		return "", errors.New("state value is nil")
	} else if  verifierCode == ""  {
		return "", errors.New("verifier code value is nil")
	}

	url := u.oauth.AuthCodeURL(state, oauth2.AccessTypeOffline,
		 oauth2.S256ChallengeOption(verifierCode))

	if url == "" {
		return "", errors.New("url is nil")
	}

	return url, nil
}