package sessions

import (
	"log"
	"time"

	db "github.com/eliezedeck/goweb/db/rethinkdb"
	"github.com/satori/go.uuid"

	r "github.com/dancannon/gorethink"
)

const (
	sessionIdleTimeout    = 120
	sessionKeepAlive      = 15
	sessionRegularCleanup = 60
)

// Session represents a User session, pretty much like PHP_SESSION in the old
// days. It is an abstraction to the `session` table in the database.
// We are using database-backed session in order to be able to communicate the
// Session to other hosts on the cluster for load-balancing purpose.
type Session struct {
	uuid string
	data interface{}

	keepAliveRunning bool
	keepAliveStop    chan struct{}
	selfWatchRunning bool
	selfWatchStop    chan struct{}
}

// NewSession returns a new *Session which is given a unique UUID
func NewSession(data interface{}) *Session {
	u := uuid.NewV4()
	uuid := u.String()
	s := &Session{
		uuid: uuid,
		data: data,
	}

	// Record the Session in the database
	_, err := r.Table("sessions").Insert(map[string]interface{}{
		"id":   uuid,
		"time": time.Now().Unix(),
		"data": data,
	}).RunWrite(db.S)
	if err != nil {
		log.Fatalf("Could not Insert() Session into the database: %s", err)
	}

	// Register the Session as active
	registerActiveSession <- s

	log.Printf("New Session '%s' created", uuid)
	go s.selfWatch()

	return s
}

// ResumeSession returns the Session from the cache or the database
func ResumeSession(uuid string) (*Session, error) {
	srequest := sessionRequest{
		uuid:     uuid,
		response: make(chan *Session),
	}
	s := <-srequest.response
	if s != nil {
		log.Printf("Session '%s' resumed from cache", uuid)
		s.Touch()
		return s, nil
	}

	res, err := r.Table("sessions").Get(uuid).Run(db.S)
	if err != nil {
		return nil, err
	}

	var entry map[string]interface{}
	err = res.One(&entry)
	if err != nil {
		return nil, err
	}

	s = &Session{
		uuid: entry["id"].(string),
		data: entry["data"],
	}
	s.Touch()

	// Register as Active Session from the database
	registerActiveSession <- s

	log.Printf("Session '%s' resumed from database", uuid)
	return s, nil
}

// ClearExpiredSessions deletes all expired Sessions from the database
func ClearExpiredSessions() {
	threshold := time.Now().Unix() - sessionIdleTimeout
	resp, err := r.Table("sessions").Filter(func(sess r.Term) r.Term {
		return sess.Field("time").Lt(threshold)
	}).Delete().RunWrite(db.S)
	if err != nil {
		log.Fatalln("Could not delete expired sessions:", err)
	}
	if resp.Deleted > 0 {
		log.Printf("Deleted %d expired Sessions", resp.Deleted)
	}
}

// RegularlyClearExpiredSessions starts a goroutine which calls ClearExpiredSessions()
// regularly. It initially does a clean-up (which briefly blocks) and then starts
// the regular clean-up job in the background.
func RegularlyClearExpiredSessions() {
	ClearExpiredSessions()
	go func() {
		for {
			time.Sleep(time.Second * sessionRegularCleanup)
			ClearExpiredSessions()
		}
	}()
}

// GetUUID returns the UUID of the Session
func (s *Session) GetUUID() string {
	return s.uuid
}

// Touch keeps the Session alive
func (s *Session) Touch() {
	_, err := r.Table("sessions").Get(s.uuid).Update(map[string]interface{}{
		"time": time.Now().Unix(),
	}, r.UpdateOpts{
		Durability: "soft",
	}).RunWrite(db.S)
	if err != nil {
		log.Fatalf("Could not Touch() Session '%s' error: %s", s.uuid, err)
	}
}

// KeepAlive keeps the Session alive until the returned channel is closed. Note
// that the channel needs to closed, otherwise the Session will never die.
func (s *Session) KeepAlive() chan<- struct{} {
	if s.keepAliveRunning {
		return s.keepAliveStop
	}

	s.keepAliveRunning = true
	s.keepAliveStop = make(chan struct{})
	go s.doKeepAlive()

	return s.keepAliveStop
}

func (s *Session) doKeepAlive() {
	defer func() {
		s.keepAliveRunning = false
	}()

	// Initial touch
	s.Touch()

	for {
		select {
		case <-time.After(time.Second * sessionKeepAlive):
			s.Touch()
			log.Printf("Session '%s' touched ...", s.uuid)
		case <-s.keepAliveStop:
			log.Printf("Session Keep alive for '%s' ended.", s.uuid)
			return
		}
	}
}
