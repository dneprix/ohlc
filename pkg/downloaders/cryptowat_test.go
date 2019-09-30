package downloaders

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/dneprix/ohlc/pkg/assets"
)

func TestNewCryptowatDownloader(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	expectedDB := sqlx.NewDb(mockDB, "sqlmock")
	expectedLogger, _ := test.NewNullLogger()
	expectedName := "CRYPTOWAT"
	expectedWaitTime := cryptowatWaitTime

	actual := NewCryptowatDownloader(expectedDB, expectedLogger)

	assert.Equal(t, expectedLogger, actual.logger.Logger)
	assert.Equal(t, expectedDB, actual.db)
	assert.Equal(t, expectedName, actual.name)
	assert.Equal(t, expectedWaitTime, actual.wait)
}

func TestCryptowatDownloader_DownloadCandlesSuccess(t *testing.T) {
	d := &CryptowatDownloader{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`{
          "result": {
            "60": [
              [1481634360, 782.14, 782.14, 781.13, 781.13, 1.92525],
              [1481634420, 782.02, 782.06, 781.94, 781.98, 2.37578],
              [1481634480, 781.39, 781.94, 781.15, 781.94, 1.68882]
            ]
          },
          "allowance": {
            "cost": 28071405,
            "remaining": 6843524322,
            "upgrade": "Upgrade for a higher allowance"
          }
        }
        `))
		return
	}))

	asset := &assets.Asset{
		ID:  1,
		URL: server.URL,
	}
	actual, err := d.DownloadCandles(asset)
	assert.NoError(t, err)
	assert.Len(t, actual, 3)
}

func TestCryptowatDownloader_DownloadCandlesFailHttp(t *testing.T) {
	d := &CryptowatDownloader{}

	asset := &assets.Asset{
		ID:  1,
		URL: "TEST_URL",
	}
	actual, err := d.DownloadCandles(asset)
	assert.Nil(t, actual)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Get HTTP Request fail")
}

func TestCryptowatDownloader_DownloadCandlesFailReadBody(t *testing.T) {
	d := &CryptowatDownloader{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
		return
	}))

	asset := &assets.Asset{
		ID:  1,
		URL: server.URL,
	}
	actual, err := d.DownloadCandles(asset)
	assert.Nil(t, actual)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Read response body fail")
}

func TestCryptowatDownloader_DownloadCandlesFailParse(t *testing.T) {
	d := &CryptowatDownloader{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<>"))
		return
	}))

	asset := &assets.Asset{
		ID:  1,
		URL: server.URL,
	}
	actual, err := d.DownloadCandles(asset)
	assert.Nil(t, actual)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Parse response fail")
}

func TestCryptowatResponsePeriod_UnmarshalJSON_Success(t *testing.T) {
	period := &CryptowatResponsePeriod{}
	err := period.UnmarshalJSON([]byte("[1481634360, 781.14, 782.14, 781.13, 781.12, 1.92525]"))
	assert.NoError(t, err)
	assert.Equal(t, period.CloseTime, int64(1481634360))
	assert.Equal(t, period.OpenPrice, float32(781.14))
	assert.Equal(t, period.HighPrice, float32(782.14))
	assert.Equal(t, period.LowPrice, float32(781.13))
	assert.Equal(t, period.ClosePrice, float32(781.12))
	assert.Equal(t, period.Volume, float32(1.92525))
}

func TestCryptowatResponsePeriod_UnmarshalJSON_Fail(t *testing.T) {
	period := &CryptowatResponsePeriod{}
	err := period.UnmarshalJSON([]byte("<>"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Parse response period fail")
}
