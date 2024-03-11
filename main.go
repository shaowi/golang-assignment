package main

import (
	"container/list"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type QueueBroker struct {
	queueMap   map[string]*list.List
	quitSignal chan struct{}
}

func NewQueueBroker() *QueueBroker {
	return &QueueBroker{
		queueMap:   make(map[string]*list.List),
		quitSignal: make(chan struct{}),
	}
}

func (qb *QueueBroker) Enqueue(queueName string, item interface{}) {
	if _, ok := qb.queueMap[queueName]; !ok {
		qb.queueMap[queueName] = list.New()
	}
	qb.queueMap[queueName].PushBack(item)
}

func (qb *QueueBroker) Dequeue(queueName string, timeout int) (interface{}, bool) {
	queue, ok := qb.queueMap[queueName]
	if !ok || queue.Len() == 0 {
		// Wait for a message or timeout
		select {
		case <-time.After(time.Duration(timeout) * time.Second): // Timeout reached
			if qb.queueMap[queueName] == nil || qb.queueMap[queueName].Len() == 0 {
				return nil, false
			}
		case <-qb.quitSignal: // Stop signal received
			return nil, false
		}
	}

	// Dequeue the first item
	item := queue.Remove(queue.Front())
	return item, true
}

func (qb *QueueBroker) handlePut(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path[1:]
	message := r.URL.Query().Get("v")

	if queueName == "" || message == "" {
		http.Error(w, "Queue name or message parameter is missing", http.StatusBadRequest)
		return
	}

	qb.Enqueue(queueName, message)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message '%s' added to the queue '%s'\n", message, queueName)
}

func (qb *QueueBroker) handleGet(w http.ResponseWriter, r *http.Request) {
	queueName := r.URL.Path[1:]
	timeoutStr := r.URL.Query().Get("timeout")

	if queueName == "" {
		http.Error(w, "Queue name parameter is missing", http.StatusBadRequest)
		return
	}

	if timeoutStr == "" {
		timeoutStr = "0" // Default to 0 seconds
	}

	timeout, err := time.ParseDuration(timeoutStr + "s")
	if err != nil {
		http.Error(w, "Invalid timeout duration", http.StatusBadRequest)
		return
	}

	item, ok := qb.Dequeue(queueName, int(timeout.Seconds()))
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Send the message in the response body
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, item)
}

func (qb *QueueBroker) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		qb.handleGet(w, r)
		break
	case http.MethodPut:
		qb.handlePut(w, r)
		break
	default:
		break
	}
}

func main() {
	// Defaults to 8080 port if unspecified
	port := flag.Int("port", 8080, "the port number to listen on")
	flag.Parse()

	// Create a new queue broker
	broker := NewQueueBroker()

	// Start the HTTP server
	fmt.Printf("HTTP server listening on port %d...\n", *port)
	http.HandleFunc("/", broker.handleRequest)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Println("Error:", err)
	}

	// Clean up the broker
	close(broker.quitSignal)
}
