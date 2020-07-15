package run

import (
	"context"
	"time"

	"github.com/bakito/dns-checker/pkg/check"
	log "github.com/sirupsen/logrus"
)

var workerChannel = make(chan chan work)

type collector struct {
	work chan work // receives jobs to send to workers
	end  chan bool // when receives bool stops workers
}

func startDispatcher(workerCount int) collector {
	var i int
	var workers []worker
	input := make(chan work) // channel to receive work
	end := make(chan bool)   // channel to spin down workers
	collector := collector{work: input, end: end}

	for i < workerCount {
		i++
		log.WithField("worker", i).Info("starting worker")
		worker := worker{
			id:            i,
			channel:       make(chan work),
			workerChannel: workerChannel,
			end:           make(chan bool)}
		worker.Start()
		workers = append(workers, worker) // stores worker
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop() // stop worker
				}
				return
			case work := <-input:
				worker := <-workerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			}
		}
	}()

	return collector
}

type work struct {
	ctx         context.Context
	interval    time.Duration
	resultsChan chan execution
	target      check.Address
	chk         check.Check
}

type worker struct {
	id            int
	workerChannel chan chan work // used to communicate between dispatcher and workers
	channel       chan work
	end           chan bool
}

// start worker
func (w *worker) Start() {
	go func() {
		for {
			w.workerChannel <- w.channel // when the worker is available place channel in queue
			select {
			case job := <-w.channel: // worker has received job
				runCheck(job, w.id) // do work
			case <-w.end:
				return
			}
		}
	}()
}

// end worker
func (w *worker) Stop() {
	log.WithField("worker", w.id).Info("worker is stopping")
	w.end <- true
}
