package melody

type hub struct {
	sessions   map[*Session]bool
	broadcast  chan *envelope
	register   chan *Session
	unregister chan *Session
}

func newHub() *hub {
	return &hub{
		sessions:   make(map[*Session]bool),
		broadcast:  make(chan *envelope),
		register:   make(chan *Session),
		unregister: make(chan *Session),
	}
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			h.sessions[s] = true
		case s := <-h.unregister:
			if _, ok := h.sessions[s]; ok {
				delete(h.sessions, s)
				close(s.output)
				s.conn.Close()
			}
		case m := <-h.broadcast:
			for s := range h.sessions {
				if m.filter != nil {
					if m.filter(s) {
						s.writeMessage(m)
					}
				} else {
					s.writeMessage(m)
				}
			}
		}
	}
}
