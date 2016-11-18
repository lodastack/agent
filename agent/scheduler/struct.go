package scheduler

import (
	"time"

	"github.com/lodastack/agent/agent/common"
	"github.com/lodastack/log"
)

type Scheduler struct {
	ticker    *time.Ticker
	quit      chan int
	interval  int
	collector common.Collector
	running   bool
}

func NewScheduler(interval int, collector common.Collector) *Scheduler {
	if interval <= 0 {
		// change it to 10s,  0 will panic agent
		interval = 10
	}
	scheduler := Scheduler{collector: collector, running: false}
	scheduler.interval = interval
	scheduler.ticker = time.NewTicker(time.Duration(interval) * time.Second)
	scheduler.quit = make(chan int)
	return &scheduler
}

func (self *Scheduler) setTicker(interval int) {
	if interval <= 0 {
		return
	}
	self.interval = interval
	self.ticker = time.NewTicker(time.Duration(interval) * time.Second)
}

func (self *Scheduler) stop() {
	if self.running {
		self.quit <- 1
		self.running = false
		log.Info("scheduler stopped: ", self.collector.Description())
	}
}

func (self *Scheduler) run() {
	if self.running {
		return
	}
	self.running = true
	log.Info("scheduler running: ", self.collector.Description())
	for {
		select {
		case <-self.ticker.C:
			self.collector.Run()
		case <-self.quit:
			return
		}
	}
}
