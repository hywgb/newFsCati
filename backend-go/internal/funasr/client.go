package funasr

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	url         string
	conn        *websocket.Conn
	OnTranscript func(string)
}

type asrMsg struct {
	Text   string `json:"text"`
	Result string `json:"result"`
}

func New(wsURL string) *Client { return &Client{url: wsURL} }

func (c *Client) Connect() error {
	if c.url == "" { return nil }
	u, _ := url.Parse(c.url)
	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	conn, _, err := dialer.Dial(u.String(), http.Header{"Origin": {"asr-gateway"}})
	if err != nil { return err }
	c.conn = conn
	return nil
}

func (c *Client) SendPCM(pcm []byte) error {
	if c.conn == nil { return nil }
	return c.conn.WriteMessage(websocket.BinaryMessage, pcm)
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *Client) ReadLoop() {
	if c.conn == nil { return }
	for {
		t, data, err := c.conn.ReadMessage()
		if err != nil { log.Printf("funasr read: %v", err); return }
		if t == websocket.TextMessage {
			msg := string(data)
			// try json
			var m asrMsg
			if json.Unmarshal(data, &m) == nil {
				if m.Text != "" && c.OnTranscript != nil { c.OnTranscript(m.Text) }
				if m.Result != "" && c.OnTranscript != nil { c.OnTranscript(m.Result) }
				continue
			}
			if c.OnTranscript != nil { c.OnTranscript(msg) }
		}
	}
}