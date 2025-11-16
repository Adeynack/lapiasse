package moneydance

import (
	"encoding/json/v2"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Boolean bool

func (b *Boolean) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*b = false
		return nil
	}

	*b = string(data) == "y"

	return nil
}

type UnixDate time.Time

func (d *UnixDate) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "0" {
		*d = UnixDate(time.Time{})
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// Parse the value into an int64
	unixTimestamp, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing unix timestamp %q: %w", value, err)
	}

	*d = UnixDate(time.Unix(unixTimestamp, 0))

	return nil
}

type IntegerDate time.Time

func (d *IntegerDate) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "0" {
		*d = IntegerDate(time.Time{})
		return nil
	}

	var strIntDate string
	if err := json.Unmarshal(data, &strIntDate); err != nil {
		return err
	}
	if len(strIntDate) != 8 {
		return fmt.Errorf("expecting an 8 character string: %s", strIntDate)
	}

	t, err := time.Parse("20060102", strIntDate)
	if err != nil {
		return fmt.Errorf("parsing MdIntDate %q: %w", strIntDate, err)
	}

	*d = IntegerDate(t)

	return nil
}

type TransactionStatus string

const (
	TransactionStatusUncleared  TransactionStatus = ""
	TransactionStatusCleared    TransactionStatus = "X"
	TransactionStatusReconciled TransactionStatus = "x"
)

func (s *TransactionStatus) UnmarshalJSON(data []byte) error {
	var strStatus string
	if err := json.Unmarshal(data, &strStatus); err != nil {
		return err
	}

	switch TransactionStatus(strStatus) {
	case TransactionStatusUncleared:
	case TransactionStatusCleared:
	case TransactionStatusReconciled:
	default:
		return fmt.Errorf("unknown TransactionStatus value: %q", strStatus)
	}

	*s = TransactionStatus(strStatus)

	return nil
}

type BigInt big.Int

func (bi *BigInt) UnmarshalJSON(data []byte) error {
	var strValue string
	if err := json.Unmarshal(data, &strValue); err != nil {
		return err
	}

	var bigIntValue big.Int
	if _, ok := bigIntValue.SetString(strValue, 10); !ok {
		return fmt.Errorf("parsing BigInt value: %q", strValue)
	}

	*bi = BigInt(bigIntValue)

	return nil
}
