package session

import (
	"encoding/gob"
	"log"
)

// Flash is the type that defines what
// a flash message holds. Currently, it just
// holds a message string. In the future
// it can hold message type(e.g. errorType, etc...)
type Flash struct {
	Message string
}

func init() {
	gob.Register(&Flash{}) // This allows type to be stored in session
}

// SetFlash encodes a string into the _flash sesison var.
func (s *Session) SetFlash(message string) {
	fm := Flash{Message: message}
	s.AddFlash(fm)
}

// GetFlashes returns all flash items in a session
func (s *Session) GetFlashes() []*Flash {
	var flashes []*Flash
	if data := s.Flashes(); len(data) > 0 {
		for _, flash := range data {
			val, ok := flash.(*Flash)
			if !ok {
				log.Println("could not get flash from session: not correct type")
				continue
			}
			flashes = append(flashes, val)
		}
	}
	return flashes
}
