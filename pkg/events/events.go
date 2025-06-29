package events

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Event struct {
	Timestamp string
	Data      string
	Key       string
	Tags      map[string]string
}

func NewEvent(key string, data string) Event {
	return Event{
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
		Tags:      make(map[string]string),
		Key:       key,
	}
}

type IEvents interface {
	Emit(event Event) error
}

type IMultiEvents interface {
	Emit(key string, data string) error
	Stop()
}

type MultiEvents struct {
	events     []IEvents
	bufferLen  int
	tags       map[string]string
	emitCh     chan Event
	quitEmitCh chan struct{}
	stopped    bool
}

func NewMultiEvents(bufferLen int, events ...IEvents) *MultiEvents {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		hostname = "unknown"
	}

	tags := make(map[string]string)
	tags["hostname"] = hostname

	me := &MultiEvents{
		events:     events,
		bufferLen:  bufferLen,
		tags:       make(map[string]string),
		emitCh:     make(chan Event, bufferLen),
		quitEmitCh: make(chan struct{}),
		stopped:    false,
	}

	go me.startWorker()
	return me
}

func (l *MultiEvents) Emit(key string, data string) error {
	m := NewEvent(key, data)
	l.emit(m)
	return nil
}

func (l *MultiEvents) Stop() {
	if l.stopped {
		return
	}
}

func (l *MultiEvents) processMetric(metric Event) {
	var didLog = false

	for _, logger := range l.events {

		metric.Tags = l.tags

		if err := logger.Emit(metric); err != nil {
			fallbackErrorLog(fmt.Sprintln("Error logging metric: ", err))
		} else {
			didLog = true
		}
	}

	if !didLog {
		fallbackLog(metric)
	}
}

func (l *MultiEvents) emit(m Event) {

	select {
	case l.emitCh <- m:
		break
	default:
		fallbackErrorLog("Channel overflow detected: " + m.Key)
		go func() {
			l.processMetric(m)
		}()
	}
}

func (l *MultiEvents) startWorker() {
	defer close(l.emitCh)

	for {
		select {
		case entry := <-l.emitCh:
			l.processMetric(entry)
		case <-l.quitEmitCh:
			for entry := range l.emitCh {
				l.processMetric(entry)
			}
			return
		}
	}
}

func (l *Event) formatTags() string {
	if len(l.Tags) == 0 {
		return ""
	}

	var builder strings.Builder

	for key, value := range l.Tags {
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(value)
		builder.WriteString(",")
	}

	result := builder.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	result += ";"

	return result
}

func fallbackErrorLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("%s - [FALLBACK] : %s\n", timestamp, message))
}

func fallbackLog(m Event) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("%s - [FALLBACK] [%s] : %s = %s\n", timestamp, m.formatTags(), m.Key, m.Data))
}
