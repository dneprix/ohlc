package downloaders

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

// KrakenDownloader structure
type KrakenDownloader struct {
	*downloader
}

// NewKrakenDownloader constructor
func NewKrakenDownloader(db *sqlx.DB, logger *logrus.Logger) *KrakenDownloader {
	return &KrakenDownloader{
		newDownloader(db, logger, "KRAKEN"),
	}
}

// DownloadCandles function
func (kd *KrakenDownloader) DownloadCandles(asset *assets.Asset) ([]*candles.Candle, error) {
	// TODO: Implement DownloadCandles for KrakenDownloader
	return nil, nil
}
