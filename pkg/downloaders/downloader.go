package downloaders

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

// Downloader interface
type Downloader interface {
	Queue() chan (bool)
	Stop() chan (bool)
	Logger() *logrus.Entry
	Name() string
	DB() *sqlx.DB
	CheckWaitTimer()

	DownloadCandles(*assets.Asset) ([]*candles.Candle, error)
}

type downloader struct {
	db *sqlx.DB

	name      string
	queue     chan (bool)
	stop      chan (bool)
	wait      time.Duration
	waitTimer *time.Timer
	logger    *logrus.Entry
}

func newDownloader(db *sqlx.DB, logger *logrus.Logger, name string, wait time.Duration) *downloader {
	return &downloader{
		db:        db,
		name:      name,
		queue:     make(chan (bool), 1),
		stop:      make(chan (bool)),
		wait:      wait,
		waitTimer: time.NewTimer(0),
		logger: logger.WithFields(logrus.Fields{
			"downloader": name,
		}),
	}
}

// ProcessQueue is a goroutine for processing downloader queue
func ProcessQueue(d Downloader) {
	for {
		select {
		case <-d.Queue():
			d.Logger().Debugf("Process queue: size=%d", len(d.Queue()))

			// Get assets for downloader name
			downloaderAssets, err := assets.GeеListByDownloaderName(d.DB(), d.Name())
			if err != nil {
				d.Logger().Errorf("Get DB downloader assets fail: %s", err)
				continue
			}

			// Process each downloader asset
			for _, asset := range downloaderAssets {
				assetLogger := assets.Logger(d.Logger(), asset)

				// Check and wait timer since last downloading
				d.CheckWaitTimer()
				assetLogger.Warn("Start candles downloading")

				// Download candles
				candlesData, err := d.DownloadCandles(asset)
				if err != nil {
					assetLogger.Errorf("Download candles fail: %s", err)
					continue
				}

				// Validate candles
				if len(candlesData) == 0 {
					assetLogger.Warn("Download ZERO candles data: Nothing to save")
					continue
				}

				// Save candles
				assetLogger.Debugf("Try to save downloaded candles: %d", len(candlesData))
				if err := candles.Save(d.DB(), candlesData); err != nil {
					assetLogger.Errorf("Save candles data fail: %s", err)
					continue
				}
				assetLogger.Debug("Candles data was successfull saved")
			}
		case <-d.Stop():
			d.Logger().Warn("Stop processing queue")
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

// Queue downloader channel
func (dl *downloader) Queue() chan (bool) {
	return dl.queue
}

// Logger downloader
func (dl *downloader) Logger() *logrus.Entry {
	return dl.logger
}

// Name downloader
func (dl *downloader) Name() string {
	return dl.name
}

// DB downloader
func (dl *downloader) DB() *sqlx.DB {
	return dl.db
}

// Stop downloader channel
func (dl *downloader) Stop() chan (bool) {
	return dl.stop
}

// CheckWaitTimer downloader
func (dl *downloader) CheckWaitTimer() {
	<-dl.waitTimer.C
	dl.waitTimer.Reset(dl.wait)
}
