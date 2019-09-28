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
	Logger() *logrus.Entry
	Name() string
	DB() *sqlx.DB

	DownloadCandles(*assets.Asset) ([]*candles.Candle, error)
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

			// Download and save candles for each asset
			for _, asset := range downloaderAssets {
				assetLogger := assets.Logger(d.Logger(), asset)

				candlesData, err := d.DownloadCandles(asset)
				if err != nil {
					assetLogger.Errorf("Download candles fail: %s", err)
					continue
				}

				if len(candlesData) == 0 {
					assetLogger.Warn("Download ZERO candles data: Nothing to save")
					continue
				}

				assetLogger.Debugf("Try to save downloaded candles: %d", len(candlesData))
				if err := candles.Save(d.DB(), candlesData); err != nil {
					assetLogger.Errorf("Save candles data fail: %s", err)
					continue
				}
				assetLogger.Debug("Candles data was successfull saved")
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

type downloader struct {
	db *sqlx.DB

	name   string
	queue  chan (bool)
	logger *logrus.Entry
}

func newDownloader(db *sqlx.DB, logger *logrus.Logger, name string) *downloader {
	return &downloader{
		db:    db,
		name:  name,
		queue: make(chan (bool), 1),
		logger: logger.WithFields(logrus.Fields{
			"downloader": name,
		}),
	}
}

// Queue downloader
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
