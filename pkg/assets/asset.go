package assets

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Asset structure
type Asset struct {
	ID         uint   `db:"id"`
	CoinFrom   string `db:"coin_from"`
	CoinTo     string `db:"coin_to"`
	Exchange   string `db:"exchange"`
	Downloader string `db:"downloader"`
	URL        string `db:"url"`
}

// GetListByDownloaderName from database
func GetListByDownloaderName(db *sqlx.DB, name string) ([]*Asset, error) {
	assets := []*Asset{}
	if err := db.Select(&assets, "SELECT * FROM assets WHERE downloader=$1", name); err != nil {
		return nil, err
	}
	return assets, nil
}

// Logger with asset field
func Logger(logger *logrus.Entry, asset *Asset) *logrus.Entry {
	return logger.WithFields(
		logrus.Fields{
			"asset": fmt.Sprintf("%s/%s/%s", asset.CoinFrom, asset.CoinTo, asset.Exchange),
		},
	)
}
