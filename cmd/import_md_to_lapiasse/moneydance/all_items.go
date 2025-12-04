package moneydance

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
)

type AllItems struct {
	IgnoredTypes map[string]uint // Count of ignored items by type

	CurrencyById map[string]*Currency // Currencies by ID
	Accounts     []*Account           // Accounts by ID
}

func (ai *AllItems) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	ai.IgnoredTypes = make(map[string]uint)
	ai.CurrencyById = make(map[string]*Currency)
	ai.Accounts = make([]*Account, 0)

	// Read the start of the array
	token, err := dec.ReadToken()
	if errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("reading token from JSON decoder: %w", err)
	} else if token.Kind() != jsontext.BeginArray.Kind() {
		return fmt.Errorf("expected start of array token, got: %v", token)
	}

	// Look through the array elements
	for dec.PeekKind() != jsontext.EndArray.Kind() {
		rawValue, err := dec.ReadValue()
		if err != nil {
			return fmt.Errorf("reading value from JSON decoder: %w", err)
		}

		var typedItem struct {
			ObjType string `json:"obj_type"`
		}
		if err := json.Unmarshal(rawValue, &typedItem); err != nil {
			return fmt.Errorf("decoding MdAllItems item type: %w", err)
		}

		switch typedItem.ObjType {
		case "curr":
			err = ai.unmarshalCurrency(rawValue)
		case "acct":
			err = ai.unmarshalAccount(rawValue)
		default:
			ai.IgnoredTypes[typedItem.ObjType] = ai.IgnoredTypes[typedItem.ObjType] + 1
		}
		if err != nil {
			return fmt.Errorf("decoding MdAllItems item of type %q: %w", typedItem.ObjType, err)
		}
	}

	// Read the end of the array
	if _, err := dec.ReadToken(); err != nil {
		return fmt.Errorf("reading end of MdAllItems item object: %w", err)
	}

	return nil
}

func (ai *AllItems) unmarshalCurrency(rawValue jsontext.Value) error {
	var currency Currency
	if err := json.Unmarshal(rawValue, &currency); err != nil {
		return fmt.Errorf("decoding Currency: %w", err)
	}

	ai.CurrencyById[currency.Id] = &currency

	return nil
}

func (ai *AllItems) unmarshalAccount(rawValue jsontext.Value) error {
	var account Account
	if err := json.Unmarshal(rawValue, &account); err != nil {
		return fmt.Errorf("decoding Account: %w", err)
	}
	if account.Id == "" {
		return fmt.Errorf("account has empty ID: %+v", account)
	}

	ai.Accounts = append(ai.Accounts, &account)

	return nil
}
