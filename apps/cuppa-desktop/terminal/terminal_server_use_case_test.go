package terminal

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/coder/websocket"
)

func start(t *testing.T) *Server {
	t.Helper()
	s, err := Start(func() tea.Model { return quitter{} })
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func dial(ctx context.Context, url, origin string) (*websocket.Conn, error) {
	header := http.Header{}
	if origin != "" {
		header.Set("Origin", origin)
	}
	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{HTTPHeader: header})
	return conn, err
}

func TestTheWindowCanConnectWithTheSecretAddress(t *testing.T) {
	s := start(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dial(ctx, s.URL(), "http://wails.localhost")
	if err != nil {
		t.Fatalf("the window should connect: %v", err)
	}
	defer func() { _ = conn.CloseNow() }()
	_, frame, err := conn.Read(ctx)
	if err != nil || !strings.Contains(string(frame), `"type":"hello"`) {
		t.Fatalf("expected a hello frame, got %q (%v)", frame, err)
	}
}

func TestTheAddressIsLoopbackAndCarriesASecret(t *testing.T) {
	s := start(t)
	if !strings.HasPrefix(s.URL(), "ws://127.0.0.1:") {
		t.Fatalf("not a loopback address: %s", s.URL())
	}
	if parts := strings.Split(s.URL(), "/term/"); len(parts) != 2 || len(parts[1]) != 32 {
		t.Fatalf("expected a 32 character secret in the path: %s", s.URL())
	}
	other := start(t)
	if other.URL() == s.URL() {
		t.Fatal("each run needs its own address and secret")
	}
}

func TestWithoutTheSecretThereIsNothingToConnectTo(t *testing.T) {
	s := start(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wrong := s.URL()[:strings.LastIndex(s.URL(), "/")+1] + strings.Repeat("0", 32)
	if conn, err := dial(ctx, wrong, "http://wails.localhost"); err == nil {
		_ = conn.CloseNow()
		t.Fatal("a wrong secret must not connect")
	}
	bare := s.URL()[:strings.Index(s.URL(), "/term/")] + "/term/"
	if conn, err := dial(ctx, bare, "http://wails.localhost"); err == nil {
		_ = conn.CloseNow()
		t.Fatal("no secret must not connect")
	}
}

func TestAPageFromAnotherSiteCannotConnectEvenWithTheSecret(t *testing.T) {
	s := start(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if conn, err := dial(ctx, s.URL(), "https://evil.example"); err == nil {
		_ = conn.CloseNow()
		t.Fatal("a foreign origin must be refused")
	}
}

func TestHostMessagesReachTheProgramOnlyWhileItRuns(t *testing.T) {
	s := start(t)
	if s.Send(tea.KeyPressMsg{Code: 'x', Text: "x"}) {
		t.Fatal("no browser has connected, so no program is running")
	}
}
