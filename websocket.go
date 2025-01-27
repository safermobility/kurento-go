package kurento

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    interface{}
}

const ConnectionLost = -1

// Response represents server response
type Response struct {
	Jsonrpc string
	Id      int64
	Result  map[string]interface{}
	Error   *Error
}

func (r *Response) unmarshal(raw json.RawMessage) bool {
	if err := json.Unmarshal(raw, r); err != nil {
		return false
	}
	return r.Id > 0 && (r.Result != nil || r.Error != nil)
}

type Event struct {
	Jsonrpc string
	Method  string
	Params  map[string]interface{}
	Error   *Error
}

func (e *Event) unmarshal(raw json.RawMessage) bool {
	if err := json.Unmarshal(raw, e); err != nil {
		return false
	}
	return e.Method == "onEvent"
}

type Connection struct {
	clientId  *atomic.Int64
	eventId   float64
	clients   threadsafeClientMap
	host      string
	ws        *websocket.Conn
	SessionId string
	events    threadsafeSubscriberMap
	eChan     chan Event
	Dead      chan bool
	IsDead    bool
}

type threadsafeClientMap struct {
	clients map[int64]chan Response
	lock    sync.RWMutex
}

type threadsafeSubscriberMap struct {
	subscribers map[string]map[string]map[string]eventHandler // eventName -> objectId -> handlerId -> handler.
	lock        sync.RWMutex
}

func NewConnection(host string) (*Connection, error) {
	c := new(Connection)

	c.clientId = &atomic.Int64{}
	c.events = threadsafeSubscriberMap{
		subscribers: make(map[string]map[string]map[string]eventHandler),
	}
	c.eChan = make(chan Event, 1)
	c.clients = threadsafeClientMap{
		clients: make(map[int64]chan Response),
	}
	c.Dead = make(chan bool, 1)

	var err error
	conf, err := websocket.NewConfig(host+"/kurento", "http://127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("kurento: error creating new config: %v", err)
	}
	conf.Dialer = &net.Dialer{Timeout: 5 * time.Second}
	c.ws, err = websocket.DialConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("kurento: error dialing: %w", err)
	}
	c.host = host
	go c.handleEvents()
	go c.handleResponse()
	return c, nil
}

func (c *Connection) Create(m IMediaObject, options map[string]interface{}) error {
	elem := &MediaObject{}
	elem.setConnection(c)
	return elem.Create(m, options)
}

func (c *Connection) Close() error {
	return c.ws.Close()
}

func (c *Connection) handleResponse() {
	for { // run forever
		if logLevel >= LogLevelSilly {
			log.Println("KURENTO WS: starting receive loop")
		}
		var incoming json.RawMessage
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			log.Printf("Error receiving on websocket %s", err)
			c.IsDead = true
			c.Dead <- true
			close(c.eChan)
			break
		}
		if logLevel >= LogLevelSilly {
			log.Println("KURENTO WS: received message")
		}

		if logLevel > LogLevelDebug {
			log.Printf("RAW %s", message)
		}

		err = json.Unmarshal([]byte(message), &incoming)
		if err != nil {
			log.Printf("Invalid JSON received on WebSocket: %s", err)
			continue
		}

		r := Response{}
		ev := Event{}
		isResponse := r.unmarshal(incoming)
		isEvent := ev.unmarshal(incoming)

		if isResponse {
			// If sessionId has been set/changed, save the new one
			if sessionID, ok := r.Result["sessionId"].(string); ok && sessionID != "" && c.SessionId != sessionID {
				if logLevel > 0 {
					log.Println("sessionId returned ", r.Result["sessionId"])
				}
				c.SessionId = sessionID
			}
			if logLevel > 0 {
				log.Printf("Response: %v", r)
			}
			// if websocket client exists, send response to the channel
			if c.clients.clients[r.Id] != nil {
				if logLevel >= LogLevelSilly {
					log.Printf("KURENTO WS: sending response to client: %d\n", r.Id)
				}
				c.clients.clients[r.Id] <- r
				// chanel is read, we can delete it
				close(c.clients.clients[r.Id])
				delete(c.clients.clients, r.Id)
				if logLevel >= LogLevelSilly {
					log.Printf("KURENTO WS: finished handling response in client: %d\n", r.Id)
				}
			} else if logLevel > 0 {
				log.Println("Dropped message because there is no client ", r.Id)
				log.Println(r)
			}
		} else if isEvent {
			if logLevel >= LogLevelSilly {
				log.Printf("KURENTO WS: sending event to client: %s\n", ev.Method)
			}
			c.eChan <- ev
			if logLevel >= LogLevelSilly {
				log.Printf("KURENTO WS: finished handling event in client: %s\n", ev.Method)
			}
		} else {
			log.Println("Unsupported message from KMS: ", message)
		}

	}
}

func (c *Connection) handleEvents() {
	for ev := range c.eChan { // run until the channel is closed
		val := ev.Params["value"].(map[string]interface{})
		if logLevel > 0 {
			log.Printf("Received event value %v", val)
		}

		t := val["type"].(string)
		objectId := val["object"].(string)

		data := val["data"].(map[string]interface{})

		if handlers, ok := c.events.subscribers[t]; ok {
			if logLevel >= LogLevelSilly {
				log.Printf("KURENTO WS: start event loop 0: %s\n", t)
			}
			if objHandlers, ok := handlers[objectId]; ok {
				if logLevel >= LogLevelSilly {
					log.Printf("KURENTO WS: start event loop 1: %s - %s\n", t, objectId)
				}
				for name, handler := range objHandlers {
					if logLevel >= LogLevelSilly {
						log.Printf("KURENTO WS: start event loop 2: %s - %s - %s\n", t, objectId, name)
					}
					handler(data)
					if logLevel >= LogLevelSilly {
						log.Printf("KURENTO WS:   end event loop 2: %s - %s - %s\n", t, objectId, name)
					}
				}
				if logLevel >= LogLevelSilly {
					log.Printf("KURENTO WS:   end event loop 1: %s - %s\n", t, objectId)
				}
			}
			if logLevel >= LogLevelSilly {
				log.Printf("KURENTO WS:   end event loop 0: %s\n", t)
			}
		}
	}
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
	if c.IsDead {
		errchan := make(chan Response, 1)
		errresp := Response{
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}
		if idInterface, ok := req["id"]; ok {
			if id, ok := idInterface.(int64); ok {
				errresp.Id = id
			}
		}
		errchan <- errresp
		return errchan
	}

	reqId := c.clientId.Add(1)
	req["id"] = reqId
	if c.SessionId != "" {
		req["sessionId"] = c.SessionId
	}
	c.clients.clients[reqId] = make(chan Response)
	if logLevel > 0 {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	err := websocket.JSON.Send(c.ws, req)
	if err != nil {
		log.Printf("Error sending on websocket %s", err)
		c.Dead <- true
		c.IsDead = true

		delete(c.clients.clients, reqId)

		errchan := make(chan Response, 1)
		errresp := Response{
			Id: req["id"].(int64),
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}

		errchan <- errresp
		return errchan
	}
	return c.clients.clients[reqId]
}

func (c *Connection) Subscribe(event, objectId, handlerId string, handler eventHandler) {
	var oh map[string]map[string]eventHandler
	var ok bool

	if oh, ok = c.events.subscribers[event]; !ok {
		c.events.subscribers[event] = make(map[string]map[string]eventHandler)
		oh = c.events.subscribers[event]
	}

	var he map[string]eventHandler
	if he, ok = oh[objectId]; !ok {
		oh[objectId] = make(map[string]eventHandler)
		he = oh[objectId]
	}

	he[handlerId] = handler
}

func (c *Connection) Unsubscribe(event, objectId, handlerId string) {
	var oh map[string]map[string]eventHandler
	var he map[string]eventHandler
	var ok bool

	if oh, ok = c.events.subscribers[event]; !ok {
		return // not found
	}

	if he, ok = oh[objectId]; !ok {
		return // not found
	}

	delete(he, handlerId)
}
