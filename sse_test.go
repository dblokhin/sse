package sse

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	// result is used to capture each yield call.
	type result struct {
		ev  Event
		err error
	}

	tests := []struct {
		name      string
		input     string
		stopAfter int // if > 0, yield returns false after that many calls (simulate early termination)
		want      []result
	}{
		{
			name:  "Simple event with id",
			input: "id: 1\n\n",
			want: []result{
				{ev: Event{ID: []byte("1")}},
			},
		},
		{
			name:  "Multiple data lines",
			input: "data: first line\ndata: second line\n\n",
			want: []result{
				{ev: Event{Data: []byte("first line\nsecond line")}},
			},
		},
		{
			name:  "Multiple fields",
			input: "id: 1\nevent: message\ndata: hello\nretry: 5000\n\n",
			want: []result{
				{ev: Event{
					ID:    []byte("1"),
					Event: []byte("message"),
					Data:  []byte("hello"),
					Retry: []byte("5000"),
				}},
			},
		},
		{
			name:  "Skip comment",
			input: ": this is a comment\ndata: hello\n\n",
			want: []result{
				{ev: Event{Data: []byte("hello")}},
			},
		},
		{
			name:  "Invalid sequence",
			input: "invalid\ndata: hello\n\n",
			want: []result{
				{
					err: fmt.Errorf("%w: %q", ErrInvalidSequence, "invalid"),
				},
				{
					ev: Event{Data: []byte("hello")},
				},
			},
		},
		{
			name:  "Empty field skipped",
			input: "id: \ndata: hello\n\n",
			want: []result{
				{ev: Event{Data: []byte("hello")}},
			},
		},
		{
			name:  "Last event on EOF",
			input: "id: 2\ndata: bye",
			want: []result{
				{ev: Event{ID: []byte("2"), Data: []byte("bye")}},
			},
		},
		{
			name:      "Stop early",
			input:     "id: 1\n\nid: 2\n\n",
			stopAfter: 1,
			want: []result{
				{ev: Event{ID: []byte("1")}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var results []result
			count := 0
			// yield collects each event/error. It simulates early termination if stopAfter is set.
			yield := func(ev Event, err error) bool {
				results = append(results, result{ev: ev, err: err})
				count++
				if tc.stopAfter > 0 && count == tc.stopAfter {
					return false
				}
				return true
			}

			// Call the iterator returned by Read
			r := strings.NewReader(tc.input)
			Read(r)(yield)

			if len(results) != len(tc.want) {
				t.Fatalf("expected %d results, got %d", len(tc.want), len(results))
			}
			// Compare each result
			for i, got := range results {
				want := tc.want[i]
				if !bytes.Equal(got.ev.ID, want.ev.ID) {
					t.Errorf("result %d: expected ID %q, got %q", i, want.ev.ID, got.ev.ID)
				}
				if !bytes.Equal(got.ev.Event, want.ev.Event) {
					t.Errorf("result %d: expected Event %q, got %q", i, want.ev.Event, got.ev.Event)
				}
				if !bytes.Equal(got.ev.Data, want.ev.Data) {
					t.Errorf("result %d: expected Data %q, got %q", i, want.ev.Data, got.ev.Data)
				}
				if !bytes.Equal(got.ev.Retry, want.ev.Retry) {
					t.Errorf("result %d: expected Retry %q, got %q", i, want.ev.Retry, got.ev.Retry)
				}
				if (got.err == nil) != (want.err == nil) {
					t.Errorf("result %d: expected error %v, got %v", i, want.err, got.err)
				} else if got.err != nil && want.err != nil {
					if got.err.Error() != want.err.Error() {
						t.Errorf("result %d: expected error %q, got %q", i, want.err.Error(), got.err.Error())
					}
				}
			}
		})
	}
}
