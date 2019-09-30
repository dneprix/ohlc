package downloaders

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

const krakenWaitTime = 10 * time.Second

// KrakenDownloader structure
type KrakenDownloader struct {
	*downloader
}

// NewKrakenDownloader constructor
func NewKrakenDownloader(db *sqlx.DB, logger *logrus.Logger) *KrakenDownloader {
	return &KrakenDownloader{
		newDownloader(db, logger, "KRAKEN", krakenWaitTime),
	}
}

// DownloadCandles function
func (kd *KrakenDownloader) DownloadCandles(asset *assets.Asset) ([]*candles.Candle, error) {
	// TODO: Implement DownloadCandles for KrakenDownloader
	return nil, nil
}
