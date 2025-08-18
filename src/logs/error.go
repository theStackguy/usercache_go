package logs

import "errors"

var (
	ErrGuid          = errors.New("issue occurred during guid creation")
	ErrUser          = errors.New("user not found")
	ErrPayLoad       = errors.New("id/Key/Value seems missing")
	ErrCacheUpdate   = errors.New("issue occured while adding/updating the data")
	ErrUserExpired   = errors.New("requested User is expired")
	ErrCacheExpired  = errors.New("requested cache is expired")
	ErrExpired       = errors.New("expired data couldnt be updated")
	ErrReadUser      = errors.New("requested User is not available")
	ErrReadUserToken = errors.New("requested User token is not available")
	ErrReadCacheKey  = errors.New("requested cache key couldnt be found")
)
