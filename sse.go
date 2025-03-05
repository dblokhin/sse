package sse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Event is message unit in SSE
type Event struct {
	ID    []byte // The event ID to set the EventSource object's last event ID value.
	Event []byte // A string identifying the type of event described.
	Data  []byte // The data field for the message.
	Retry []byte // The reconnection time (milliseconds). If a non-integer value is specified, the field is ignored.
}

var (
	splitField = []byte(`:`)
	idField    = []byte(`id`)
	eventField = []byte(`event`)
	dataField  = []byte(`data`)
	retryField = []byte(`retry`)

	ErrInvalidSequence = errors.New("invalid event sequence")
)

// Read returns iterator over streaming in input reader.
// It automatically skips comments. Returns ErrInvalidSequence with the problem line if got invalid line sequence.
func Read(input io.Reader) func(yield func(Event, error) bool) {
	return func(yield func(Event, error) bool) {
		scan := bufio.NewScanner(input)
		event := Event{}
		updated := false

		for scan.Scan() {
			payload := scan.Bytes()

			if len(payload) == 0 {
				// fire event
				if updated && !yield(event, nil) {
					return
				}
				event = Event{}
				updated = false
				continue
			}

			// parse line
			del := bytes.IndexByte(payload, ':')
			if del == 0 {
				continue // skip comment
			}
			if del < 0 {
				err := fmt.Errorf("%w: %q", ErrInvalidSequence, string(payload))
				if !yield(Event{}, err) {
					return
				}
				continue // skip invalid sequence
			}

			field, content := payload[:del], bytes.TrimSpace(payload[del+1:])

			// update fields
			updated = true
			switch {
			case bytes.EqualFold(field, idField):
				event.ID = content
			case bytes.EqualFold(field, eventField):
				event.Event = content
			case bytes.EqualFold(field, dataField):
				if event.Data == nil {
					event.Data = make([]byte, 0)
				} else {
					event.Data = append(event.Data, '\n')
				}
				event.Data = append(event.Data, content...)
			case bytes.EqualFold(field, retryField):
				event.Retry = content
			}
		}

		// fire last event if exists
		if updated && !yield(event, nil) {
			return
		}

		if err := scan.Err(); err != nil {
			yield(Event{}, err)
		}
	}
}
