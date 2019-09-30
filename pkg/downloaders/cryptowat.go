package downloaders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/dneprix/ohlc/pkg/assets"
	"github.com/dneprix/ohlc/pkg/candles"
)

const cryptowatWaitTime = 20 * time.Second

// CryptowatDownloader structure
type CryptowatDownloader struct {
	*downloader
}

// NewCryptowatDownloader constructor
func NewCryptowatDownloader(db *sqlx.DB, logger *logrus.Logger) *CryptowatDownloader {
	return &CryptowatDownloader{
		newDownloader(db, logger, "CRYPTOWAT", cryptowatWaitTime),
	}
}

// DownloadCandles function
func (cd *CryptowatDownloader) DownloadCandles(asset *assets.Asset) ([]*candles.Candle, error) {

	res, err := http.Get(asset.URL)
	if err != nil {
		return nil, fmt.Errorf("Get HTTP Request fail: %s", err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Read response body fail: %s", err)
	}

	candlesResponse := CryptowatResponse{}
	if err := json.Unmarshal(data, &candlesResponse); err != nil {
		return nil, fmt.Errorf("Parse response fail: %s", err)
	}

	candlesData := make([]*candles.Candle, 0, len(candlesResponse.Result.Period))
	for _, period := range candlesResponse.Result.Period {
		candlesData = append(candlesData, &candles.Candle{
			AssetID:    asset.ID,
			CloseTime:  period.CloseTime,
			OpenPrice:  period.OpenPrice,
			HighPrice:  period.HighPrice,
			LowPrice:   period.LowPrice,
			ClosePrice: period.ClosePrice,
			Volume:     period.Volume,
		})
	}

	return candlesData, nil
}

// CryptowatResponse structure
type CryptowatResponse struct {
	Result struct {
		Period []CryptowatResponsePeriod `json:"60"`
	} `json:"result"`
}

// CryptowatResponsePeriod structure
type CryptowatResponsePeriod struct {
	CloseTime  int64
	OpenPrice  float32
	HighPrice  float32
	LowPrice   float32
	ClosePrice float32
	Volume     float32
}

// UnmarshalJSON for CryptowatResponsePeriod
func (crp *CryptowatResponsePeriod) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&crp.CloseTime,
		&crp.OpenPrice,
		&crp.HighPrice,
		&crp.LowPrice,
		&crp.ClosePrice,
		&crp.Volume,
	}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return fmt.Errorf("Parse response period fail: %s", err)
	}
	return nil
}
