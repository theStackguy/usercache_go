package src

import (
	"time"
)

const (
	DefaultSessionTokenExpiry = defaultSessionTokenTime * time.Minute
	DefaultRefreshTokenExpiry = defaultRefreshTokenTime * time.Hour
)

func (session *Session)generateSessionRefreshToken(sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration) (string, string, error) {
	sessionToken, err := generateToken(session_token_length)
	if err == nil {
		refreshToken, err := generateToken(refresh_token_length)
		if err == nil {
			if sessionTokenExpiryTime <= 0 {
				session.SessionExpiry = time.Now().Add(DefaultSessionTokenExpiry)
			} else {
				session
			}
			if refreshTokenExpiryTime <= 0 {
				session.RefreshExpiry = time.Now().Add(DefaultRefreshTokenExpiry)
			}
			return sessionToken, refreshToken, nil
		}
	}
	return "", "", err

}

func (s *Session) checkTokenExpired() {
	if s.SessionExpiry.IsZero() {
		s.Mu.Lock()
		s.SessionToken = ""
		s.Err = errZeroExpiry
		s.Mu.Unlock()
		return
	}
	if time.Now().After(s.SessionExpiry) {
		if !time.Now().After(s.RefreshExpiry) {
			sessionToken, err := generateToken(session_token_length)
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
