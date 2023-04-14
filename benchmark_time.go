package main

import (
	"net/http"
	"time"
)

type BenchmarkTime struct {
	bench         *Benchmark
	executionTime time.Duration
}

func NewBenchmarkTime(context *Context) *BenchmarkTime {
	collector := make(chan *Record, 1e7)
	return &BenchmarkTime{
		executionTime: context.config.executionWindow,
		bench:         &Benchmark{context, collector}}
}

func (b *BenchmarkTime) Run() {
	jobs := make(chan *http.Request, b.bench.c.config.concurrency*GoMaxProcs)

	for i := 0; i < b.bench.c.config.concurrency; i++ {
		go NewHTTPWorker(b.bench.c, jobs, b.bench.collector).Run()
	}

	ticker := time.NewTicker(b.executionTime)
	base, _ := NewHTTPRequest(b.bench.c.config)

	for {
		select {
		case <-ticker.C:
			close(jobs)
			<-b.bench.c.stop
			break
		default:
			jobs <- CopyHTTPRequest(b.bench.c.config, base)
		}
	}

}
