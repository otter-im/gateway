package model

import (
	"github.com/go-oauth2/oauth2/v4"
	"time"
)

type Token struct {
	tableName struct{} `json:"-" pg:"auth_tokens"`

	ClientID    string `json:"client_id" pg:"client_id"`
	UserID      string `json:"user_id" pg:"user_id"`
	RedirectURI string `json:"redirect_uri" pg:"redirect_uri"`
	Scope       string `json:"scope" pg:"scope"`

	Code                string    `json:"code" pg:"code"`
	CodeChallenge       string    `json:"code_challenge" pg:"code_challenge"`
	CodeChallengeMethod string    `json:"code_challenge_method" pg:"code_challenge_method"`
	CodeCreateAt        time.Time `json:"code_create_at" pg:"code_create_at"`
	CodeExpiresAt       time.Time `json:"code_expires_at" pg:"code_expires_at"`

	Access          string    `json:"access" pg:"access"`
	AccessCreateAt  time.Time `json:"access_create_at" pg:"access_create_at"`
	AccessExpiresAt time.Time `json:"access_expires_at" pg:"access_expires_at"`

	Refresh          string    `json:"refresh" pg:"refresh"`
	RefreshCreateAt  time.Time `json:"refresh_create_at" pg:"refresh_create_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at" pg:"refresh_expires_at"`
}

func (t *Token) New() oauth2.TokenInfo {
	return new(Token)
}

func (t *Token) GetClientID() string {
	return t.ClientID
}

// SetClientID the client id
func (t *Token) SetClientID(clientID string) {
	t.ClientID = clientID
}

// GetUserID the user id
func (t *Token) GetUserID() string {
	return t.UserID
}

// SetUserID the user id
func (t *Token) SetUserID(userID string) {
	t.UserID = userID
}

// GetRedirectURI redirect URI
func (t *Token) GetRedirectURI() string {
	return t.RedirectURI
}

// SetRedirectURI redirect URI
func (t *Token) SetRedirectURI(redirectURI string) {
	t.RedirectURI = redirectURI
}

// GetScope get scope of authorization
func (t *Token) GetScope() string {
	return t.Scope
}

// SetScope get scope of authorization
func (t *Token) SetScope(scope string) {
	t.Scope = scope
}

// GetCode authorization code
func (t *Token) GetCode() string {
	return t.Code
}

// SetCode authorization code
func (t *Token) SetCode(code string) {
	t.Code = code
}

// GetCodeCreateAt create Time
func (t *Token) GetCodeCreateAt() time.Time {
	return t.CodeCreateAt
}

// SetCodeCreateAt create Time
func (t *Token) SetCodeCreateAt(createAt time.Time) {
	t.CodeCreateAt = createAt
}

// GetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *Token) GetCodeExpiresIn() time.Duration {
	return t.CodeExpiresAt.Sub(time.Now())
}

// SetCodeExpiresIn the lifetime in seconds of the authorization code
func (t *Token) SetCodeExpiresIn(exp time.Duration) {
	t.CodeExpiresAt = time.Now().Add(exp)
}

// GetCodeChallenge challenge code
func (t *Token) GetCodeChallenge() string {
	return t.CodeChallenge
}

// SetCodeChallenge challenge code
func (t *Token) SetCodeChallenge(code string) {
	t.CodeChallenge = code
}

// GetCodeChallengeMethod challenge method
func (t *Token) GetCodeChallengeMethod() oauth2.CodeChallengeMethod {
	return oauth2.CodeChallengeMethod(t.CodeChallengeMethod)
}

// SetCodeChallengeMethod challenge method
func (t *Token) SetCodeChallengeMethod(method oauth2.CodeChallengeMethod) {
	t.CodeChallengeMethod = string(method)
}

// GetAccess access Token
func (t *Token) GetAccess() string {
	return t.Access
}

// SetAccess access Token
func (t *Token) SetAccess(access string) {
	t.Access = access
}

// GetAccessCreateAt create Time
func (t *Token) GetAccessCreateAt() time.Time {
	return t.AccessCreateAt
}

// SetAccessCreateAt create Time
func (t *Token) SetAccessCreateAt(createAt time.Time) {
	t.AccessCreateAt = createAt
}

// GetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) GetAccessExpiresIn() time.Duration {
	return t.AccessExpiresAt.Sub(time.Now())
}

// SetAccessExpiresIn the lifetime in seconds of the access token
func (t *Token) SetAccessExpiresIn(exp time.Duration) {
	t.AccessExpiresAt = time.Now().Add(exp)
}

// GetRefresh refresh Token
func (t *Token) GetRefresh() string {
	return t.Refresh
}

// SetRefresh refresh Token
func (t *Token) SetRefresh(refresh string) {
	t.Refresh = refresh
}

// GetRefreshCreateAt create Time
func (t *Token) GetRefreshCreateAt() time.Time {
	return t.RefreshCreateAt
}

// SetRefreshCreateAt create Time
func (t *Token) SetRefreshCreateAt(createAt time.Time) {
	t.RefreshCreateAt = createAt
}

// GetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *Token) GetRefreshExpiresIn() time.Duration {
	return t.RefreshExpiresAt.Sub(time.Now())
}

// SetRefreshExpiresIn the lifetime in seconds of the refresh token
func (t *Token) SetRefreshExpiresIn(exp time.Duration) {
	t.RefreshExpiresAt = time.Now().Add(exp)
}
