package session

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1 //63
	letterIdxMax  = 63 / letterIdxBits   //10
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var randNum = rand.NewSource(time.Now().UnixNano())

// Session wraps gorilla's session type to add
// convenience methods to keep controller funcs cleaner
type Session struct {
	*sessions.Session
}

// GetSession gets a function from the provided store
func GetSession(store *sessions.CookieStore, r *http.Request, sessionName string) (*Session, error) {
	s, err := store.Get(r, sessionName)
	return &Session{s}, err
}

type sessionKeys struct {
	newAuthentication []byte
	newEncryption     []byte
	oldAuthentication []byte
	oldEncryption     []byte
}

func newSessionKeys() *sessionKeys {
	return &sessionKeys{
		newAuthentication: generateRandomKey(32),
		newEncryption:     generateRandomKey(32),
		oldAuthentication: generateRandomKey(32),
		oldEncryption:     generateRandomKey(32),
	}
}

func (sk *sessionKeys) rotateKeys() {
	sk.oldAuthentication = sk.newAuthentication
	sk.oldEncryption = sk.newEncryption
	sk.newAuthentication = generateRandomKey(32)
	sk.newEncryption = generateRandomKey(32)
}

// CookieStore holds all the things needed to
// manage a cookieStore
type CookieStore struct {
	Store       *sessions.CookieStore
	sks         *sessionKeys
	SessionName string
}

// NewCookieStore makes a new cookie store. <- wins the
// award of the year for the most boring comment ever.
func NewCookieStore(sessionName string) *CookieStore {
	sks := newSessionKeys()
	s := sessions.NewCookieStore(
		sks.newAuthentication,
		sks.newEncryption,
		sks.oldAuthentication,
		sks.oldEncryption,
	)
	return &CookieStore{
		Store:       s,
		sks:         sks,
		SessionName: sessionName,
	}
}

// RotateSessionKeys can be called to rotate cookie store
// keys for better security.
func (ss CookieStore) RotateSessionKeys() {
	ss.sks.rotateKeys()
	ss.Store = sessions.NewCookieStore(
		ss.sks.newAuthentication,
		ss.sks.newEncryption,
		ss.sks.oldAuthentication,
		ss.sks.oldEncryption,
	)
}

// generateRandomKey returns a random string of len n
func generateRandomKey(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, randNum.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randNum.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}
