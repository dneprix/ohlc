package candles

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSaveSuccess(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO candles")
	mock.ExpectExec("INSERT INTO candles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	candles := []*Candle{
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
	}
	err := Save(sqlxDB, candles)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveExecFail(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	expectedError := fmt.Errorf("DB error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO candles")
	mock.ExpectExec("INSERT INTO candles").WillReturnError(expectedError)

	candles := []*Candle{
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
	}
	err := Save(sqlxDB, candles)

	assert.Error(t, err)
	assert.EqualError(t, err, "Tx stmt exec fail: "+expectedError.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveCommitFail(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	expectedError := fmt.Errorf("DB error")
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO candles")
	mock.ExpectExec("INSERT INTO candles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(expectedError)

	candles := []*Candle{
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
	}
	err := Save(sqlxDB, candles)

	assert.Error(t, err)
	assert.EqualError(t, err, "Tx commit fail: "+expectedError.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
