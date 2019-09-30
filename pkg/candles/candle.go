package candles

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Candle structure
type Candle struct {
	ID         uint    `db:"id"`
	AssetID    uint    `db:"asset_id"`
	Period     uint    `db:"period"`
	CloseTime  int64   `db:"close_time"`
	OpenPrice  float32 `db:"open_price"`
	HighPrice  float32 `db:"high_price"`
	LowPrice   float32 `db:"low_price"`
	ClosePrice float32 `db:"close_price"`
	Volume     float32 `db:"volume"`
}

// Save to database
func Save(db *sqlx.DB, candles []*Candle) error {
	sqls := `INSERT INTO candles(
        asset_id,
        period,
        close_time,
        open_price,
        high_price,
        low_price,
        close_price,
        volume
      ) VALUES(
        :asset_id,
        :period,
        to_timestamp(:close_time),
        :open_price,
        :high_price,
        :low_price,
        :close_price,
        :volume
      )
      ON CONFLICT DO NOTHING;
      `

	tx := db.MustBegin()
	stmt, _ := tx.PrepareNamed(sqls)
	for _, candle := range candles {
		if _, err := stmt.Exec(candle); err != nil {
			return fmt.Errorf("Tx stmt exec fail: %s", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Tx commit fail: %s", err)
	}

	return nil
}
