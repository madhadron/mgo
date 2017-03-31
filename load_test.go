package mgo

import (
	"math/rand"
	"mgo/bson"
	"testing"
)

// Generate inserts and upserts to do.
// 50% inserts/50% upserts. Messages are
type TestMessage struct {
	upsert bool
	id     bson.ObjectId
	v      []byte
}

func randomBytes(n int) []byte {

	return v
}

func generate(ch chan *TestMessage) {
	oidPool := make([]bson.ObjectId, 1000)
	for i, _ := range oidPool {
		oidPool[i] = bson.NewObjectId()
	}

	i := 0
	for {
		var n int
		if i < 3 {
			n = 10000
		} else {
			n = 1000
		}
		v := make([]byte, n)
		rand.Read(v)

		var msg *TestMessage
		if i%2 == 0 {
			msg = &TestMessage{
				upsert: true,
				id:     oidPool[rand.Intn(len(oidPool))],
				v:      v,
			}
		} else {
			msg = &TestMessage{
				upsert: false,
				id:     bson.NewObjectId(),
				v:      v,
			}
		}

		i = (i + 1) % 100
		ch <- msg
	}
}

func push(ch chan *TestMessage) {
	// Dial

	for msg := range ch {
		// Insert or upsert
	}
}

func TestLoad(t *testing.T) {
	ch := make(chan *TestMessage, 100)
	go generate(ch)
	go push(ch)
}
