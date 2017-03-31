package mgo

import (
	"fmt"
	"math/rand"
	"mgo/bson"
	"sync/atomic"
	"testing"
	"time"
)

// Generate inserts and upserts to do.
// 50% inserts/50% upserts. Messages are
type TestMessage struct {
	upsert bool
	id     bson.ObjectId
	v      []byte
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

var metrics struct {
	inserts uint64
	upserts uint64
}

func startMetrics() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var lastInserts uint64
	var lastUpserts uint64
	for ts := range ticker.C {
		inserts := atomic.LoadUint64(&metrics.inserts)
		upserts := atomic.LoadUint64(&metrics.upserts)

		dInserts := inserts - lastInserts
		dUpserts := upserts - lastUpserts

		fmt.Printf("%v\tinserts=%d\tupserts=%d\n", ts, dInserts, dUpserts)
		lastInserts = inserts
		lastUpserts = upserts
	}
}
