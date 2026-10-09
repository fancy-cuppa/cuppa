// Package terminal serves the Cuppa terminal app to the window's web view: one
// shared Bubble Tea program behind a WebSocket that only this computer can
// reach, and only with this run's secret in the address.
package terminal

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/meta-tui/treactui/packages/tty-go/session"
	"github.com/meta-tui/treactui/packages/tty-go/socket"
)

// pageOrigins are the hosts the Wails web view loads its page from: Windows
// uses http://wails.localhost, macOS and Linux use wails://wails, and
// `wails dev` serves the page from localhost.
var pageOrigins = []string{"wails.localhost", "wails", "localhost:*", "127.0.0.1:*"}

// Server is the running terminal endpoint.
type Server struct {
	shared *session.SharedProgram
	http   *http.Server
	url    string
}

// Start listens on a free loopback port and serves the program built by
// newModel. The program is shared: reloading the window finds it as it was.
func Start(newModel func() tea.Model) (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	token, err := newToken()
	if err != nil {
		_ = listener.Close()
		return nil, err
	}
	shared := session.NewSharedProgram(newModel)
	path := "/term/" + token
	mux := http.NewServeMux()
	mux.Handle(path, socket.SharedHandler(shared, socket.Options{AllowedOrigins: pageOrigins}))
	s := &Server{
		shared: shared,
		http:   &http.Server{Handler: mux},
		url:    "ws://" + listener.Addr().String() + path,
	}
	go func() { _ = s.http.Serve(listener) }()
	return s, nil
}

// URL is the WebSocket address the page connects to.
func (s *Server) URL() string { return s.url }

// Send delivers a message to the running program and reports whether one was running.
func (s *Server) Send(msg tea.Msg) bool { return s.shared.Send(msg) }

// Close stops serving.
func (s *Server) Close() error { return s.http.Close() }

// newToken is a random path segment: the endpoint runs a program, so another
// program on this computer must not be able to guess where it is.
func newToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("no randomness available: " + err.Error())
	}
	return hex.EncodeToString(b), nil
}
