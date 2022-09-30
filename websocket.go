package kurento

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    string
}

const ConnectionLost = -1

// Response represents server response
type Response struct {
	Jsonrpc string
	Id      float64
	Result  map[string]interface{}
	Error   *Error
}

type Event struct {
	Jsonrpc string
	Method  string
	Params  map[string]interface{}
	Error   *Error
}

type Connection struct {
	clientId  float64
	eventId   float64
	clients   threadsafeClientMap
	host      string
	ws        *websocket.Conn
	SessionId string
	events    threadsafeSubscriberMap
	Dead      chan bool
	IsDead    bool
}

type threadsafeClientMap struct {
	clients map[float64]chan Response
	lock    sync.RWMutex
}

type threadsafeSubscriberMap struct {
	subscribers map[string]map[string]map[string]eventHandler // eventName -> objectId -> handlerId -> handler.
	lock        sync.RWMutex
}

func NewConnection(host string) (*Connection, error) {
	c := new(Connection)

	c.events = threadsafeSubscriberMap{
		subscribers: make(map[string]map[string]map[string]eventHandler),
	}
	c.clients = threadsafeClientMap{
		clients: make(map[float64]chan Response),
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
		r := Response{}
		ev := Event{}
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			log.Printf("Error receiving on websocket %s", err)
			c.IsDead = true
			c.Dead <- true
			break
		}

		if debug {
			log.Printf("RAW %s", message)
		}

		// Decode into both possible types. One should be valid
		_ = json.Unmarshal([]byte(message), &r)
		_ = json.Unmarshal([]byte(message), &ev)

		isResponse := r.Id > 0 && r.Result != nil
		isEvent := ev.Method == "onEvent"

		//websocket.JSON.Receive(c.ws, &r)
		if isResponse {
			// If sessionId has been set/changed, save the new one
			if sessionID, ok := r.Result["sessionId"].(string); ok && sessionID != "" && c.SessionId != sessionID {
				if debug {
					log.Println("sessionId returned ", r.Result["sessionId"])
				}
				c.SessionId = sessionID
			}
			if debug {
				log.Printf("Response: %v", r)
			}
			// if webscocket client exists, send response to the channel
			c.clients.lock.RLock()
			c.events.lock.RLock()
			defer c.clients.lock.RUnlock()
			defer c.events.lock.RUnlock()
			if c.clients.clients[r.Id] != nil {
				c.clients.clients[r.Id] <- r
				// chanel is read, we can delete it
				c.clients.lock.Lock()
				defer c.clients.lock.Unlock()
				close(c.clients.clients[r.Id])
				delete(c.clients.clients, r.Id)
			} else if debug {
				log.Println("Dropped message because there is no client ", r.Id)
				log.Println(r)
			}
		} else if isEvent {

			val := ev.Params["value"].(map[string]interface{})
			if debug {
				log.Printf("Received event value %v", val)
			}

			t := val["type"].(string)
			objectId := val["object"].(string)

			data := val["data"].(map[string]interface{})

			if handlers, ok := c.events.subscribers[t]; ok {
				if objHandlers, ok := handlers[objectId]; ok {
					for _, handler := range objHandlers {
						handler(data)
					}
				}
			}
		} else {
			log.Println("Unsupported message from KMS: ", message)
		}

	}
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
	c.clients.lock.Lock()
	defer c.clients.lock.Unlock()

	if c.IsDead {
		errchan := make(chan Response, 1)
		errresp := Response{
			Id: req["id"].(float64),
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}
		errchan <- errresp
		return errchan
	}

	c.clientId++
	req["id"] = c.clientId
	if c.SessionId != "" {
		req["sessionId"] = c.SessionId
	}
	c.clients.clients[c.clientId] = make(chan Response)
	if debug {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	err := websocket.JSON.Send(c.ws, req)
	if err != nil {
		log.Printf("Error sending on websocket %s", err)
		c.Dead <- true
		c.IsDead = true

		delete(c.clients.clients, c.clientId)

		errchan := make(chan Response, 1)
		errresp := Response{
			Id: req["id"].(float64),
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}

		errchan <- errresp
		return errchan
	}
	return c.clients.clients[c.clientId]
}

func (c *Connection) Subscribe(event, objectId, handlerId string, handler eventHandler) {
	var oh map[string]map[string]eventHandler
	var ok bool

	c.events.lock.Lock()
	defer c.events.lock.Unlock()

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

	c.events.lock.Lock()
	defer c.events.lock.Unlock()

	if oh, ok = c.events.subscribers[event]; !ok {
		return // not found
	}

	if he, ok = oh[objectId]; !ok {
		return // not found
	}

	delete(he, handlerId)
}