package src

import (
	"time"
)

const (
	DefaultSessionTokenExpiry = DefaultSessionTokenTime * time.Minute
	DefaultRefreshTokenExpiry = DefaultRefreshTokenTime * time.Hour
)

func (session *Session) generateSessionRefreshToken(sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration) {
	sessionToken, sessiontokenerr := generateToken(session_token_length)
	if sessiontokenerr != nil {
		session.Err = errSessionTokenGen
		return
	}
	session.SessionToken = sessionToken

	refreshToken, refreshtokenerr := generateToken(refresh_token_length)
	if refreshtokenerr != nil {
		session.Err = errRefershTokenGen
		return
	}
	session.RefreshToken = refreshToken

	if sessionTokenExpiryTime <= 0 {
		session.SessionExpiry = time.Now().Add(DefaultSessionTokenExpiry)
	} else {
		session.SessionExpiry = time.Now().Add(sessionTokenExpiryTime)
	}
	if refreshTokenExpiryTime <= 0 {
		session.RefreshExpiry = time.Now().Add(DefaultRefreshTokenExpiry)
	} else {
		session.RefreshExpiry = time.Now().Add(refreshTokenExpiryTime)
	}
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


func ( u *User)verifySessionCredentials(sessionid string, sessiontoken string) error {
	
     u.Mu.RLock()
	 session,exist := u.Sessions[sessionid]
	 if exist {
         if  (u.CurrentSessionId == sessionid) {
           
		 }
	 } 
		return  errSession
	 
}
