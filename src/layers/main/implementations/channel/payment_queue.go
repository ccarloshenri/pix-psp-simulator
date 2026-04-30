package channelqueue

import "pix-psp-simulator/src/layers/main/interfaces"

// Queue is a buffered channel-backed implementation of interfaces.PaymentQueue.
// The buffer size is set at construction time to prevent enqueuers from
// blocking under normal load.
type Queue struct {
	ch chan interfaces.PaymentJob
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{ch: make(chan interfaces.PaymentJob, bufferSize)}
}

func (q *Queue) Enqueue(job interfaces.PaymentJob) { q.ch <- job }
func (q *Queue) Jobs() <-chan interfaces.PaymentJob { return q.ch }
