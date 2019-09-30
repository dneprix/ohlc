package schedulers

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/downloaders"
)

var timerDuration = time.Minute

// Scheduler structure
type Scheduler struct {
	downloaders []downloaders.Downloader
	logger      *logrus.Logger
	stop        chan (bool)
}

// NewSheduler constructor
func NewSheduler(logger *logrus.Logger) *Scheduler {
	return &Scheduler{
		logger: logger,
		stop:   make(chan (bool)),
	}
}

// Add downloader to scheduler
func (s *Scheduler) Add(d downloaders.Downloader) {
	s.downloaders = append(s.downloaders, d)
}

// Run scheduler
func (s *Scheduler) Run() {
	s.logger.Info("Run Scheduler")

	// Start processing downloaders queues
	for _, d := range s.downloaders {
		go downloaders.ProcessQueue(d)
	}

	// Push each downloader queue by timer
	// Only 1 slot available for waiting in queue. Skip if queue is full
	for range time.Tick(timerDuration) {
		s.logger.Info("Run All Downloaders")
		for _, d := range s.downloaders {
			select {
			case d.Queue() <- true:
				d.Logger().Debugf("Add to downloader queue: size=%d", len(d.Queue()))
			case <-s.stop:
				d.Logger().Warn("Stop scheduler")
				return
			default:
				d.Logger().Warnf("Skip adding. Downloader queue is full: size=%d", len(d.Queue()))
			}
		}
	}
}
