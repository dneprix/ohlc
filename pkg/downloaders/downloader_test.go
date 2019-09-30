package downloaders

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

type mockDownloader struct {
	*downloader
	TestDownloadCandles func() ([]*candles.Candle, error)
}

func (m *mockDownloader) DownloadCandles(asset *assets.Asset) ([]*candles.Candle, error) {
	return m.TestDownloadCandles()
}

func Test_newDownloader(t *testing.T) {
	expectedLogger, _ := test.NewNullLogger()

	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	expectedDB := sqlx.NewDb(mockDB, "sqlmock")

	expectedName := "DOWNLOADER"
	expectedWaitTime := 10 * time.Millisecond

	actual := newDownloader(expectedDB, expectedLogger, expectedName, expectedWaitTime)

	assert.Equal(t, expectedLogger, actual.logger.Logger)
	assert.Equal(t, expectedDB, actual.db)
	assert.Equal(t, expectedName, actual.name)
	assert.Equal(t, expectedWaitTime, actual.wait)
}

func TestProcessQueue(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	name := "TEST_DOWNLOADER"
	logger, _ := test.NewNullLogger()

	d := &mockDownloader{
		downloader: &downloader{
			db:        db,
			name:      name,
			queue:     make(chan (bool), 1),
			stop:      make(chan (bool)),
			wait:      time.Millisecond,
			waitTimer: time.NewTimer(0),
			logger: logger.WithFields(logrus.Fields{
				"downloader": name,
			}),
		},
	}

	go func() {
		<-time.After(time.Millisecond)
		close(d.Stop())
	}()
	d.Queue() <- true
	ProcessQueue(d)
}

func TestProcessDownloaderSuccess(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	name := "TEST_DOWNLOADER"
	logger, _ := test.NewNullLogger()

	d := &mockDownloader{
		downloader: &downloader{
			db:        db,
			name:      name,
			queue:     make(chan (bool), 1),
			stop:      make(chan (bool)),
			wait:      time.Millisecond,
			waitTimer: time.NewTimer(0),
			logger: logger.WithFields(logrus.Fields{
				"downloader": name,
			}),
		},
		TestDownloadCandles: func() ([]*candles.Candle, error) {
			return []*candles.Candle{
				{
					AssetID:    1,
					Period:     60,
					CloseTime:  1569563400,
					OpenPrice:  7937.2001953125,
					HighPrice:  7937.2001953125,
					LowPrice:   7937.60009765625,
					ClosePrice: 7937.60009765625,
					Volume:     0.008394920267164707,
				},
			}, nil
		},
	}

	asset := &assets.Asset{
		ID:         1,
		CoinFrom:   "TEST_COIN_FROM",
		CoinTo:     "TEST_COIN_TO",
		Exchange:   "TEST_EXCHANGE",
		Downloader: "TEST_DOWNLOADER",
		URL:        "TEST_URL",
	}

	mock.ExpectQuery(
		"SELECT \\* FROM assets WHERE downloader=\\$1",
	).WithArgs(asset.Downloader).WillReturnRows(
		sqlmock.NewRows([]string{"id", "coin_from", "coin_to", "exchange", "downloader", "url"}).
			AddRow(
				asset.ID,
				asset.CoinFrom,
				asset.CoinTo,
				asset.Exchange,
				asset.Downloader,
				asset.URL,
			))
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO candles")
	mock.ExpectExec("INSERT INTO candles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ProcessDownloader(d)
}

func TestProcessDownloaderFailDownloadCandles(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	name := "TEST_DOWNLOADER"
	logger, _ := test.NewNullLogger()

	d := &mockDownloader{
		downloader: &downloader{
			db:        db,
			name:      name,
			queue:     make(chan (bool), 1),
			stop:      make(chan (bool)),
			wait:      time.Millisecond,
			waitTimer: time.NewTimer(0),
			logger: logger.WithFields(logrus.Fields{
				"downloader": name,
			}),
		},
		TestDownloadCandles: func() ([]*candles.Candle, error) {
			return nil, fmt.Errorf("DB error")
		},
	}

	asset := &assets.Asset{
		ID:         1,
		CoinFrom:   "TEST_COIN_FROM",
		CoinTo:     "TEST_COIN_TO",
		Exchange:   "TEST_EXCHANGE",
		Downloader: "TEST_DOWNLOADER",
		URL:        "TEST_URL",
	}

	mock.ExpectQuery(
		"SELECT \\* FROM assets WHERE downloader=\\$1",
	).WithArgs(asset.Downloader).WillReturnRows(
		sqlmock.NewRows([]string{"id", "coin_from", "coin_to", "exchange", "downloader", "url"}).
			AddRow(
				asset.ID,
				asset.CoinFrom,
				asset.CoinTo,
				asset.Exchange,
				asset.Downloader,
				asset.URL,
			))

	ProcessDownloader(d)
}

func TestProcessDownloaderFailNoCandles(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	name := "TEST_DOWNLOADER"
	logger, _ := test.NewNullLogger()

	d := &mockDownloader{
		downloader: &downloader{
			db:        db,
			name:      name,
			queue:     make(chan (bool), 1),
			stop:      make(chan (bool)),
			wait:      time.Millisecond,
			waitTimer: time.NewTimer(0),
			logger: logger.WithFields(logrus.Fields{
				"downloader": name,
			}),
		},
		TestDownloadCandles: func() ([]*candles.Candle, error) {
			return []*candles.Candle{}, nil
		},
	}

	asset := &assets.Asset{
		ID:         1,
		CoinFrom:   "TEST_COIN_FROM",
		CoinTo:     "TEST_COIN_TO",
		Exchange:   "TEST_EXCHANGE",
		Downloader: "TEST_DOWNLOADER",
		URL:        "TEST_URL",
	}

	mock.ExpectQuery(
		"SELECT \\* FROM assets WHERE downloader=\\$1",
	).WithArgs(asset.Downloader).WillReturnRows(
		sqlmock.NewRows([]string{"id", "coin_from", "coin_to", "exchange", "downloader", "url"}).
			AddRow(
				asset.ID,
				asset.CoinFrom,
				asset.CoinTo,
				asset.Exchange,
				asset.Downloader,
				asset.URL,
			))

	ProcessDownloader(d)
}

func TestProcessDownloaderFailSaveCandles(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	name := "TEST_DOWNLOADER"
	logger, _ := test.NewNullLogger()

	d := &mockDownloader{
		downloader: &downloader{
			db:        db,
			name:      name,
			queue:     make(chan (bool), 1),
			stop:      make(chan (bool)),
			wait:      time.Millisecond,
			waitTimer: time.NewTimer(0),
			logger: logger.WithFields(logrus.Fields{
				"downloader": name,
			}),
		},
		TestDownloadCandles: func() ([]*candles.Candle, error) {
			return []*candles.Candle{
				{
					AssetID:    1,
					Period:     60,
					CloseTime:  1569563400,
					OpenPrice:  7937.2001953125,
					HighPrice:  7937.2001953125,
					LowPrice:   7937.60009765625,
					ClosePrice: 7937.60009765625,
					Volume:     0.008394920267164707,
				},
			}, nil
		},
	}

	asset := &assets.Asset{
		ID:         1,
		CoinFrom:   "TEST_COIN_FROM",
		CoinTo:     "TEST_COIN_TO",
		Exchange:   "TEST_EXCHANGE",
		Downloader: "TEST_DOWNLOADER",
		URL:        "TEST_URL",
	}

	mock.ExpectQuery(
		"SELECT \\* FROM assets WHERE downloader=\\$1",
	).WithArgs(asset.Downloader).WillReturnRows(
		sqlmock.NewRows([]string{"id", "coin_from", "coin_to", "exchange", "downloader", "url"}).
			AddRow(
				asset.ID,
				asset.CoinFrom,
				asset.CoinTo,
				asset.Exchange,
				asset.Downloader,
				asset.URL,
			))
	expectedError := fmt.Errorf("DB error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO candles")
	mock.ExpectExec("INSERT INTO candles").WillReturnError(expectedError)

	ProcessDownloader(d)
}

func Test_downloader_Queue(t *testing.T) {
	expected := make(chan (bool))
	dl := &downloader{
		queue: expected,
	}
	actual := dl.Queue()
	assert.Equal(t, expected, actual)
}

func Test_downloader_Logger(t *testing.T) {
	logger, _ := test.NewNullLogger()
	expected := logrus.NewEntry(logger)
	dl := &downloader{
		logger: expected,
	}
	actual := dl.Logger()
	assert.Equal(t, expected, actual)
}

func Test_downloader_Name(t *testing.T) {
	expected := "DOWNLOADER"
	dl := &downloader{
		name: expected,
	}
	actual := dl.Name()
	assert.Equal(t, expected, actual)
}

func Test_downloader_DB(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	expected := sqlx.NewDb(mockDB, "sqlmock")

	dl := &downloader{
		db: expected,
	}
	actual := dl.DB()
	assert.Equal(t, expected, actual)
}

func Test_downloader_Stop(t *testing.T) {
	expected := make(chan (bool))
	dl := &downloader{
		stop: expected,
	}
	actual := dl.Stop()
	assert.Equal(t, expected, actual)
}

func Test_downloader_CheckWaitTimer(t *testing.T) {
	expected := time.Duration(time.Millisecond)
	dl := &downloader{
		wait:      expected,
		waitTimer: time.NewTimer(0),
	}

	// Check zero wait time
	startTime := time.Now()
	dl.CheckWaitTimer()
	endTime := time.Now()
	assert.Equal(t, time.Duration(0), endTime.Sub(startTime).Truncate(time.Millisecond))

	// Check wait time
	startTime = time.Now()
	dl.CheckWaitTimer()
	endTime = time.Now()
	assert.Equal(t, expected, endTime.Sub(startTime).Truncate(time.Millisecond))

	// Check wait time + expected process
	startTime = time.Now()
	time.Sleep(expected)
	dl.CheckWaitTimer()
	endTime = time.Now()
	assert.Equal(t, expected, endTime.Sub(startTime).Truncate(time.Millisecond))
}
