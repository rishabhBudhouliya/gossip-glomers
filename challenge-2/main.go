package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Counter struct {
	mu sync.Mutex
	id string
}

func (c *Counter) generateId() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.id = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
		time.Now().Second(), time.Now().Nanosecond()+rand.IntN(100), time.UTC).String()
}

func NewCounter() *Counter {
	c := &Counter{}
	c.generateId()
	return c
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		counter := NewCounter()
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "generate_ok"
		body["in_reply_to"] = body["msg_id"]
		body["id"] = counter.id
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
