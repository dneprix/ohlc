package schedulers

import (
	"testing"
	"time"

	"github.com/dneprix/ohlc/pkg/downloaders"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

type mockDownloader struct {
	downloaders.Downloader
	queue  chan (bool)
	stop   chan (bool)
	logger *logrus.Entry
}

func (m *mockDownloader) Queue() chan (bool) {
	return m.queue
}

func (m *mockDownloader) Logger() *logrus.Entry {
	return m.logger
}

func (m *mockDownloader) Stop() chan (bool) {
	return m.stop
}

func TestNewSheduler(t *testing.T) {
	logger, _ := test.NewNullLogger()
	s := NewSheduler(logger)
	assert.Len(t, s.downloaders, 0)
	assert.NotNil(t, s.logger)
	assert.Equal(t, s.logger, logger)
}

func TestScheduler_Add(t *testing.T) {
	s := &Scheduler{}
	assert.Len(t, s.downloaders, 0)
	s.Add(&mockDownloader{})
	assert.Len(t, s.downloaders, 1)
	s.Add(&mockDownloader{})
	assert.Len(t, s.downloaders, 2)
}

func TestScheduler_Run(t *testing.T) {
	logger, _ := test.NewNullLogger()
	s := &Scheduler{
		logger: logger,
		stop:   make(chan (bool)),
		downloaders: []downloaders.Downloader{
			&mockDownloader{
				queue:  make(chan (bool), 1),
				logger: logrus.NewEntry(logger),
				stop:   make(chan (bool)),
			},
		},
	}

	close(s.downloaders[0].Stop())
	go func() {
		<-time.After(10 * time.Millisecond)
		close(s.stop)
	}()
	timerDuration = time.Millisecond
	s.Run()
}
