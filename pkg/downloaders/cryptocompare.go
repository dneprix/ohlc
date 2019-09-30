package downloaders

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

const cryptocompareWaitTime = 10 * time.Second

// CryptocompareDownloader structure
type CryptocompareDownloader struct {
	*downloader
}

// NewCryptocompareDownloader constructor
func NewCryptocompareDownloader(db *sqlx.DB, logger *logrus.Logger) *CryptocompareDownloader {
	return &CryptocompareDownloader{
		newDownloader(db, logger, "CRYPTOCOMPARE", cryptocompareWaitTime),
	}
}

// DownloadCandles function
func (kd *CryptocompareDownloader) DownloadCandles(asset *assets.Asset) ([]*candles.Candle, error) {
	// TODO: Implement DownloadCandles for CryptocompareDownloader
	return nil, nil
}
