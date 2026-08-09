package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const sessionTTL = 15 * time.Minute

type Auth struct {
	DB     *database.Database
	secret []byte
}

func New(db *database.Database, jwtSecret string) *Auth {
	return &Auth{DB: db, secret: []byte(jwtSecret)}
}

// IssueSession creates the session in the database, returns the short-lived
// JWT for the client to store in sessionStorage and sends the persist cookie.
func (a *Auth) IssueSession(ctx context.Context, w http.ResponseWriter, userID uuid.UUID) (string, error) {
	sid, err := a.DB.CreateSession(ctx, userID)
	if err != nil {
		return "", err
	}

	u, err := a.DB.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":        u.ID.String(),
		"sid":        sid,
		"email":      u.Email,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"exp":        time.Now().Add(sessionTTL).Unix(),
	}
	signed, err := a.sign(claims)
	if err != nil {
		return "", err
	}

	setPersistCookie(w, sid)
	return signed, nil
}

// RefreshSession rotates the session: deletes the old row, creates a new one.
func (a *Auth) RefreshSession(ctx context.Context, w http.ResponseWriter, sessionID string) (string, error) {
	userID, err := a.DB.SessionUserID(ctx, sessionID)
	if err != nil {
		return "", err
	}
	if err := a.DB.DeleteSession(ctx, sessionID); err != nil {
		return "", err
	}
	return a.IssueSession(ctx, w, userID)
}

func (a *Auth) sign(claims jwt.MapClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(a.secret)
}

func userFromClaims(claims jwt.MapClaims) (*User, error) {
	sub, _ := claims["sub"].(string)
	id, err := uuid.Parse(sub)
	if err != nil {
		return nil, errors.New("invalid subject claim")
	}
	email, _ := claims["email"].(string)
	firstName, _ := claims["first_name"].(string)
	lastName, _ := claims["last_name"].(string)
	return &User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}