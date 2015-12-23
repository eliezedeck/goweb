package sessions

import (
	"log"

	r "github.com/dancannon/gorethink"
	db "github.com/eliezedeck/goweb/db/rethinkdb"
)

type sessionRequest struct {
	uuid     string
	response chan *Session
}

var (
	activeSessions        = make(map[string]*Session)
	registerActiveSession = make(chan *Session, 2)
	getSession            = make(chan sessionRequest)
	dropSession           = make(chan *Session, 2)
)

// Start sets-up and starts the necessary ceremony to support Session
func Start() {
	go startKeeper()

	log.Println("Sessions keeper started...")
}

func startKeeper() {
	for {
		select {
		case getting := <-getSession:
			if s, ok := activeSessions[getting.uuid]; ok {
				getting.response <- s
			} else {
				getting.response <- nil
			}
		case session := <-registerActiveSession:
			activeSessions[session.uuid] = session
			go session.selfWatch()
		case session := <-dropSession:
			delete(activeSessions, session.uuid)
		}
	}
}

func (s *Session) selfWatch() {
	if s.selfWatchRunning {
		return
	}
	s.selfWatchRunning = true
	s.selfWatchStop = make(chan struct{})

	cur, err := r.Table("sessions").Get(s.uuid).Changes(r.ChangesOpts{
		IncludeInitial: true,
	}).Pluck(map[string]interface{}{
		"new_val": map[string]interface{}{
			"id": true,
		},
	}).Run(db.S)
	if err != nil {
		log.Fatalln("Could not Session.selfWatch():", err)
	}

	// Ensure that cur.Next() that will be called later on in an endless loop will
	// be able to stop by getting cur.Close() called when the s.selfWatchStop is
	// signaled
	closeCursorInDefer := true
	cursorClosedInDefer := make(chan struct{})
	defer func() {
		if closeCursorInDefer {
			cur.Close()
			close(cursorClosedInDefer)
		}
		s.selfWatchRunning = false
	}()
	go func() {
		select {
		case <-s.selfWatchStop:
			closeCursorInDefer = false
			cur.Close()
		case <-cursorClosedInDefer:
			return
		}
	}()

	log.Printf("Self watch started for '%s' ...", s.uuid)

	var entry map[string]interface{}
	for cur.Next(&entry) {
		// log.Println("Entry:", entry)
		if entry["new_val"] == nil {
			log.Printf("Session '%s' has been deleted", s.uuid)

			// Delete from activeSessions
			dropSession <- s

			break
		}
	}

	log.Printf("Self watch ended for '%s' ...", s.uuid)
}
