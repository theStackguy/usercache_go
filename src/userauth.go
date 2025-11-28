package src

import (
	"time"
)

const (
	DefaultSessionTokenExpiry = DefaultSessionTokenTime * time.Minute
	DefaultRefreshTokenExpiry = DefaultRefreshTokenTime * time.Hour
)

func generateSessionRefreshToken(sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration) (string, string, error) {
	sessionToken, err := generateToken(SESSION_TOKEN_LENGTH)
	if err == nil {
		refreshToken, err := generateToken(REFRESH_TOKEN_LENGTH)
		if err == nil {
			if sessionTokenExpiryTime <= ZERO {
				sessionTokenExpiryTime = DefaultSessionTokenExpiry
			}
			if refreshTokenExpiryTime <= ZERO {
				refreshTokenExpiryTime = DefaultRefreshTokenExpiry
			}
			return sessionToken, refreshToken, nil
		}
	}
	return "", "", err

}

func (s *Session) checkTokenExpired() {
	if s.Expiry.IsZero() {
		s.Mu.Lock()
		s.SessionToken = ""
		s.Err = errZeroExpiry
		s.Mu.Unlock()
		return
	}
	if time.Now().After(s.Expiry) {
		if !time.Now().After(s.RefreshExpiry) {
			sessionToken, err := generateToken(SESSION_TOKEN_LENGTH)
			if err != nil {
				s.Mu.Lock()
				s.SessionToken = ""
				s.Err = errTokenGen
				s.Mu.Unlock()
				return
			}
			s.Mu.Lock()
			s.SessionToken = sessionToken
			s.Err = nil
			s.Mu.Unlock()
			return
		}
		s.Mu.Lock()
		s.SessionToken = ""
		s.Err = errAuth
		s.Mu.Unlock()
		return 
	}
		
}
