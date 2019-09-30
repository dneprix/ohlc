package assets

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGeеListByDownloaderNameSuccess(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	expectedAssets := []*Asset{
		{
			ID:         1,
			CoinFrom:   "TEST_COIN_FROM",
			CoinTo:     "TEST_COIN_TO",
			Exchange:   "TEST_EXCHANGE",
			Downloader: "TEST_DOWNLOADER",
			URL:        "TEST_URL",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "coin_from", "coin_to", "exchange", "downloader", "url"}).
		AddRow(
			expectedAssets[0].ID,
			expectedAssets[0].CoinFrom,
			expectedAssets[0].CoinTo,
			expectedAssets[0].Exchange,
			expectedAssets[0].Downloader,
			expectedAssets[0].URL,
		)
	mock.ExpectQuery(
		"SELECT \\* FROM assets WHERE downloader=\\$1",
	).WithArgs(expectedAssets[0].Downloader).WillReturnRows(rows)

	actualAssets, err := GeеListByDownloaderName(sqlxDB, expectedAssets[0].Downloader)

	assert.Equal(t, expectedAssets, actualAssets)
	assert.NoError(t, err)
}

func TestGeеListByDownloaderNameFail(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	expectedError := fmt.Errorf("DB error")
	expectedDownloader := "DOWNLOADER"
	mock.ExpectQuery("SELECT").WithArgs(expectedDownloader).WillReturnError(expectedError)

	actualAssets, err := GeеListByDownloaderName(sqlxDB, expectedDownloader)

	assert.Nil(t, actualAssets)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedError.Error())
}

func TestLogger(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	asset := &Asset{
		ID:         1,
		CoinFrom:   "TEST_COIN_FROM",
		CoinTo:     "TEST_COIN_TO",
		Exchange:   "TEST_EXCHANGE",
		Downloader: "TEST_DOWNLOADER",
		URL:        "TEST_URL",
	}
	expectedData := "TEST_COIN_FROM/TEST_COIN_TO/TEST_EXCHANGE"

	actualLogger := Logger(logger, asset)
	assert.NotNil(t, actualLogger)
	assert.Equal(t, expectedData, actualLogger.Data["asset"])
}
