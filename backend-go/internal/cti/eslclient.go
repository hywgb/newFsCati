package cti

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Event struct {
	Headers map[string]string
	Body    string
}

type EventHandler func(ev Event)

type Client struct {
	addr     string
	password string
	conn     net.Conn
	br       *bufio.Reader
	bw       *bufio.Writer
	mu       sync.Mutex
	onEvent  EventHandler
	closed   chan struct{}
}

func NewClient(addr, password string, h EventHandler) *Client {
	return &Client{addr: addr, password: password, onEvent: h, closed: make(chan struct{})}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, 5*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	c.br = bufio.NewReader(conn)
	c.bw = bufio.NewWriter(conn)
	// read auth/request
	if _, err := c.readFrame(); err != nil {
		return err
	}
	if err := c.sendLine("auth " + c.password); err != nil {
		return err
	}
	if _, err := c.readFrame(); err != nil { // +OK accepted
		return err
	}
	// subscribe events plain ALL
	if err := c.sendLine("event plain ALL"); err != nil {
		return err
	}
	if _, err := c.readFrame(); err != nil { // +OK
		return err
	}
	go c.loop()
	return nil
}

func (c *Client) Close() error {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) loop() {
	for {
		select {
		case <-c.closed:
			return
		default:
		}
		f, err := c.readFrame()
		if err != nil {
			return
		}
		if c.onEvent != nil {
			c.onEvent(f)
		}
	}
}

func (c *Client) readFrame() (Event, error) {
	ev := Event{Headers: map[string]string{}}
	for {
		line, err := c.br.ReadString('\n')
		if err != nil {
			return ev, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			ev.Headers[key] = val
		}
	}
	if cl, ok := ev.Headers["Content-Length"]; ok {
		var n int
		fmt.Sscanf(cl, "%d", &n)
		buf := make([]byte, n)
		if _, err := ioReadFull(c.br, buf); err != nil {
			return ev, err
		}
		ev.Body = string(buf)
	}
	return ev, nil
}

func ioReadFull(r *bufio.Reader, buf []byte) (int, error) {
	n := 0
	for n < len(buf) {
		m, err := r.Read(buf[n:])
		if err != nil {
			return n, err
		}
		n += m
	}
	return n, nil
}

func (c *Client) sendLine(s string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return errors.New("not connected")
	}
	if _, err := c.bw.WriteString(s + "\n\n"); err != nil {
		return err
	}
	return c.bw.Flush()
}

func (c *Client) API(cmd string) error { return c.sendLine("api " + cmd) }
func (c *Client) BGAPI(cmd string) error { return c.sendLine("bgapi " + cmd) }

func (c *Client) UUIDKill(uuid, cause string) error {
	if cause == "" {
		cause = "NORMAL_CLEARING"
	}
	return c.API(fmt.Sprintf("uuid_kill %s %s", uuid, cause))
}

func (c *Client) UUIDAudioForkStart(uuid, wsURL, params string) error {
	p := params
	if p == "" {
		p = "{channels=1,sampling=8000,ptime=20,stream-type=sender}"
	}
	return c.BGAPI(fmt.Sprintf("uuid_audio_fork %s start %s %s", uuid, wsURL, p))
}

func (c *Client) UUIDAudioForkStop(uuid string) error {
	return c.BGAPI(fmt.Sprintf("uuid_audio_fork %s stop all", uuid))
}