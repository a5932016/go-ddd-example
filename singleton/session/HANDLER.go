package session

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Session represents a user session
type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // get current session ID
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime MaxLifeTime)
}

var (
	ErrSessionNotExisted = errors.New("session ID does not exist")
)

type SessionName string
type MaxLifeTime int64

func NewManager(provider Provider, sessionName SessionName, maxLifeTime MaxLifeTime) *Manager {
	return &Manager{
		provider:    provider,
		sessionName: sessionName,
		maxLifeTime: maxLifeTime,
	}
}

type Manager struct {
	provider    Provider
	sessionName SessionName
	maxLifeTime MaxLifeTime
}

func (manager *Manager) newSessionID() string {
	return uuid.New().String()
}

type SessionCarrier struct {
	Name    string
	ID      string
	Session Session
}

func (manager *Manager) SessionStart(unescapedID string) (SessionCarrier, error) {
	// New session
	if unescapedID == "" {
		sid := manager.newSessionID()
		session, err := manager.provider.SessionInit(sid)
		if err != nil {
			return SessionCarrier{}, errors.WithMessagef(err,
				"(sessionName: %s) provider.SessionInit(%s)", manager.sessionName, sid)
		}
		return SessionCarrier{
			Name:    string(manager.sessionName),
			ID:      url.QueryEscape(sid),
			Session: session,
		}, nil
	}

	// Read existing session
	sid, err := url.QueryUnescape(unescapedID)
	if err != nil {
		return SessionCarrier{}, errors.WithMessagef(err,
			"(SessionName: %s) url.QueryUnescape(%s)", manager.sessionName, unescapedID)
	}
	session, err := manager.provider.SessionRead(sid)
	if err != nil {
		if errors.Is(err, ErrSessionNotExisted) {
			return SessionCarrier{}, err
		}
		return SessionCarrier{}, errors.WithMessagef(err,
			"(SessionName: %s) provider.SessionRead(%s)", manager.sessionName, sid)
	}
	return SessionCarrier{
		Name:    string(manager.sessionName),
		ID:      unescapedID,
		Session: session,
	}, nil
}

func (manager *Manager) SessionDestroy(unescapedID string) error {
	if unescapedID != "" {
		sid, err := url.QueryUnescape(unescapedID)
		if err != nil {
			return errors.WithMessagef(err,
				"(SessionName: %s) url.QueryUnescape(%s)", manager.sessionName, unescapedID)
		}
		if err := manager.provider.SessionDestroy(sid); err != nil {
			return errors.WithMessagef(err,
				"(SessionName: %s) provider.SessionDestroy(%s)", manager.sessionName, sid)
		}
	}
	return nil
}

// func (manager *Manager) GC() {
// 	manager.provider.SessionGC(manager.maxLifeTime)
// 	time.AfterFunc(time.Duration(manager.maxLifeTime)*time.Second, func() {
// 		manager.GC()
// 	})
// }
