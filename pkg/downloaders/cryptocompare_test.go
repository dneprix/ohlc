package downloaders

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/dneprix/ohlc/pkg/assets"
)

func TestNewCryptocompareDownloader(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	expectedDB := sqlx.NewDb(mockDB, "sqlmock")
	expectedLogger, _ := test.NewNullLogger()
	expectedName := "CRYPTOCOMPARE"
	expectedWaitTime := cryptocompareWaitTime

	actual := NewCryptocompareDownloader(expectedDB, expectedLogger)

	assert.Equal(t, expectedLogger, actual.logger.Logger)
	assert.Equal(t, expectedDB, actual.db)
	assert.Equal(t, expectedName, actual.name)
	assert.Equal(t, expectedWaitTime, actual.wait)
}

func TestCryptocompareDownloader_DownloadCandles(t *testing.T) {
	d := &CryptocompareDownloader{}
	asset := &assets.Asset{}
	actual, err := d.DownloadCandles(asset)
	assert.Nil(t, actual)
	assert.NoError(t, err)
}
