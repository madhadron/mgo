package mgo

import (
	"fmt"
	"gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
	"math/rand"
	"os"
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
	// will need to source /opt/packetsled/etc/packetsled-pms.conf for this
	pw := os.Getenv("PMS_PACKETSLED_PASSWORD")
	un := "packetsled"
	host := "localhost"
	session, err := mgo.Dial(fmt.Sprintf("%s:%s@%s", un, pw, host))
	if err != nil {
		fmt.Println("Error openning connection: ", err)
		os.Exit(1)
	}

	// Needed to target the right db
	// will need to source /opt/packetsled/options/packetsled-ui.options for this
	envid := os.Getenv("PS_ENV_ID")
	dbName := fmt.Sprintf("probe_%s_0", envid)

	flowsColl := session.DB(dbName).C("flows")
	// eventsColl := session.DB(dbName).C("events")
	selector := map[string]interface{}{"_id": bson.NewObjectId()}
	insertDoc := map[string]interface{}{"_id": bson.NewObjectId(), "d": ""}
	upsertDoc := map[string]map[string]interface{}{"d": map[string]interface{}{"$addToSet": ""}}

	for msg := range ch {
		if msg.upsert {
			selector["_id"] = msg.id
			upsertDoc["d"]["$addToSet"] = msg.v
			atomic.AddUint64(&metrics.upserts, 1)
			flowsColl.Upsert(selector, upsertDoc)
			atomic.AddUint64(&metrics.upserts, 1)
		} else {
			insertDoc["_id"] = msg.id
			insertDoc["d"] = msg.v
			atomic.AddUint64(&metrics.inserts, 1)
			flowsColl.Insert(insertDoc)
			atomic.AddUint64(&metrics.inserts, 1)
		}
	}
}

func TestLoad(t *testing.T) {
	ch := make(chan *TestMessage, 100)
	go startMetrics()
	go generate(ch)
	go push(ch)
	time.Sleep(time.Second * 20)
}

var metrics struct {
	inserts uint64
	upserts uint64
}

func startMetrics() {
	ticker := time.NewTicker(1 * time.Second)
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
