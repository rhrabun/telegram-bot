package transcribe

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPipelineRepliesInEnqueueOrder(t *testing.T) {
	const count = 6

	var mu sync.Mutex
	var got []int
	done := make(chan struct{})

	p := newPipeline(4, 2*time.Second,
		func(_ context.Context, j *job) (string, error) {
			time.Sleep(time.Duration(count-j.messageID) * 10 * time.Millisecond)
			if j.messageID == count {
				return "", errors.New("boom")
			}
			return "text", nil
		},
		func(r result) {
			mu.Lock()
			got = append(got, r.messageID)
			if len(got) == count {
				close(done)
			}
			mu.Unlock()
		})
	p.start()

	for i := 1; i <= count; i++ {
		p.enqueue(job{messageID: i})
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for replies")
	}

	mu.Lock()
	defer mu.Unlock()
	for i, id := range got {
		if id != i+1 {
			t.Fatalf("reply order = %v, want 1..%d", got, count)
		}
	}
}
