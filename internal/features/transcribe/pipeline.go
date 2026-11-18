package transcribe

import (
	"context"
	"time"
)

type job struct {
	chatID    int64
	messageID int
	fileID    string
	fileSize  int64
}

type result struct {
	chatID    int64
	messageID int
	text      string
	err       error
}

type task struct {
	job
	res chan result
}

type workFunc func(context.Context, *job) (string, error)
type replyFunc func(result)

type pipeline struct {
	workers int
	timeout time.Duration
	work    workFunc
	reply   replyFunc

	tasks   chan task
	pending chan chan result
}

func newPipeline(workers int, timeout time.Duration, work workFunc, reply replyFunc) *pipeline {
	return &pipeline{
		workers: workers,
		timeout: timeout,
		work:    work,
		reply:   reply,
		tasks:   make(chan task, 256),
		pending: make(chan chan result, 256),
	}
}

func (p *pipeline) start() {
	for range p.workers {
		go p.worker()
	}
	go p.sender()
}

func (p *pipeline) enqueue(j job) {
	res := make(chan result, 1)
	p.pending <- res
	p.tasks <- task{job: j, res: res}
}

func (p *pipeline) worker() {
	for t := range p.tasks {
		ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
		text, err := p.work(ctx, &t.job)
		cancel()
		t.res <- result{chatID: t.chatID, messageID: t.messageID, text: text, err: err}
	}
}

func (p *pipeline) sender() {
	for res := range p.pending {
		p.reply(<-res)
	}
}
