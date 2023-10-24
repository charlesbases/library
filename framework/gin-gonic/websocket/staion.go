package websocket

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/broker"
	"github.com/charlesbases/library/framework/gin-gonic/hfwctx"
)

var es = &eventStaion{subjects: make(map[subject]subscriberGroup, 0)}

// eventStaion .
type eventStaion struct {
	client broker.Client

	subjects map[subject]subscriberGroup

	lock sync.RWMutex
}

// subscriberGroup .
type subscriberGroup map[hfwctx.ID]*subscriber

// subscriber .
type subscriber struct {
	sessionID hfwctx.ID
	subject   subject

	onEvent chan *WebSocketBroadcast
}

// InitStation .
func InitStation(client broker.Client) error {
	if client != nil {
		es.client = client
		return nil
	}
	return errors.New("broker is nil.")
}

// subscribe .
func (es *eventStaion) subscribe(subs ...*subscriber) {
	es.lock.Lock()
	for _, sub := range subs {
		if group, found := es.subjects[sub.subject]; found {
			group[sub.sessionID] = sub
		} else {
			var subscriberGroup = es.newSubscriberGroup(sub.subject)
			subscriberGroup[sub.sessionID] = sub
			es.subjects[sub.subject] = subscriberGroup
		}
	}
	es.lock.Unlock()
}

// unsubscribe .
func (es *eventStaion) unsubscribe(subs ...*subscriber) {
	es.lock.Lock()
	for _, sub := range subs {
		if group, found := es.subjects[sub.subject]; found {
			delete(group, sub.sessionID)
		}
	}
	es.lock.Unlock()
}

// newSubscriberGroup .
func (es *eventStaion) newSubscriberGroup(subject subject) subscriberGroup {
	var subscriberGroup = make(map[hfwctx.ID]*subscriber, 0)
	es.subjects[subject] = subscriberGroup

	// broker.Subscribe
	es.client.Subscribe(string(subject), es.console)

	return subscriberGroup
}

// console .
func (es *eventStaion) console(event broker.Event) error {
	var subject = subject(event.Topic())

	es.lock.RLock()
	if subscriberGroup, found := es.subjects[subject]; found {
		for _, subscriber := range subscriberGroup {
			go func() {
				mess := &WebSocketBroadcast{
					Subject: subject,
					Time:    library.NowString(),
				}
				event.Unmarshal(&mess.Data)
				subscriber.onEvent <- mess
			}()
		}
	}
	es.lock.RUnlock()

	return nil
}
