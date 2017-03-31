package mgo

import (
	"math/rand"
	"mgo/bson"
	"testing"
	"os"
	"fmt"
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
	// will need to source /opt/packetsled/etc/packetsled-pms.conf for this
	pw := os.Environ("PMS_PACKETSLED_PASSWORD")
	un := "packetsled"
	host := "localhost"
	session, err := Dial(fmt.Sprintf("%s:%s@%s", un, pw, host))
	if err := nil {
		fmt.Println("Error openning connection: ", err)
	}

	// Needed to target the right db
	// will need to source /opt/packetsled/options/packetsled-ui.options for this
	envid := os.Environ("PS_ENV_ID")
	dbName := fmt.Sprintf("probe_%s_0", envid)

	flowsColl := session.DB(dbName).C("flows")
	// eventsColl := session.DB(dbName).C("events")
	selector := map[string]interface{} {"_id": bson.NewObjectId()}
	insertDoc :=  map[string]interface{} {"_id": bson.NewObjectId(), "d": ""}
	upsertDoc :=  map[string]map[string]interface{} {"d": map[string]interface{}{"$addToSet":""}}

	for msg := range ch {
		if msg.upsert {
			selector["_id"] = msg.id
			upsertDoc["d"]["$addToSet"] = msg.v
			flowsColl.Upsert(selector, upsertDoc)
		} else {
			insertDoc["_id"] = msg.id
			insertDoc["d"] = msg.v
			flowsColl.Insert(insertDoc)
		}
	}
}

func TestLoad(t *testing.T) {
	ch := make(chan *TestMessage, 100)
	go generate(ch)
	go push(ch)
}
