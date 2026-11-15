package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// Broadcaster manages active SSE subscriptions for a specific deliberation session.
type Broadcaster struct {
	mu      sync.RWMutex
	clients map[chan protocol.EngineEvent]struct{}
	history []protocol.EngineEvent
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan protocol.EngineEvent]struct{}),
		history: make([]protocol.EngineEvent, 0, 32),
	}
}

func (b *Broadcaster) Subscribe() chan protocol.EngineEvent {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan protocol.EngineEvent, 64)
	// Replay history to newly connected client
	for _, evt := range b.history {
		ch <- evt
	}
	b.clients[ch] = struct{}{}
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan protocol.EngineEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.clients, ch)
	close(ch)
}

func (b *Broadcaster) Broadcast(evt protocol.EngineEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.history = append(b.history, evt)
	for ch := range b.clients {
		select {
		case ch <- evt:
		default:
			// Client channel full, skip to prevent blocking
		}
	}
}

// ServeHTTP handles the Server-Sent Events HTTP connection.
func (b *Broadcaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	eventChan := b.Subscribe()
	defer b.Unsubscribe(eventChan)

	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case evt, ok := <-eventChan:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, string(data))
			flusher.Flush()
		}
	}
}
