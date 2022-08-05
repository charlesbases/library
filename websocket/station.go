package websocket

import (
	"sync"

	"library"
	"library/broker"
)

var station *eventStaion

// eventStaion .
type eventStaion struct {
	event broker.Broker

	subjects map[subject]subscriberGroup

	lock sync.RWMutex
}

// subscriberGroup .
type subscriberGroup map[sessionID]*subscriber

// subscriber .
type subscriber struct {
	sessionID sessionID
	subject   subject

	onEvent chan *WebSocketBroadcast
}

// Init .
func Init(broker broker.Broker) {
	station = newStation()
	station.event = broker
}

// newStation .
func newStation() *eventStaion {
	return &eventStaion{
		subjects: make(map[subject]subscriberGroup, 0),
	}
}

// Subscribe .
func (es *eventStaion) Subscribe(subs ...*subscriber) {
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

// Unsubscribe .
func (es *eventStaion) Unsubscribe(subs ...*subscriber) {
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
	var subscriberGroup = make(map[sessionID]*subscriber, 0)
	es.subjects[subject] = subscriberGroup

	// broker.Subscribe
	es.event.Subscribe(string(subject), es.console)

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
					Time:    library.Now(),
				}
				event.Unmarshal(&mess.Data)
				subscriber.onEvent <- mess
			}()
		}
	}
	es.lock.RUnlock()

	return nil
}
