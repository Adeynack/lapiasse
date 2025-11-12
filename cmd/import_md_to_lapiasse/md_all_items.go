package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"

	"github.com/samber/lo"
)

type MdAllItems struct {
	IgnoredTypes map[string]uint
}

func (ai *MdAllItems) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	ai.IgnoredTypes = make(map[string]uint)

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
		case "c":
			// todo: handle polymorphic decoding (see tmp/polymorphic_unmarshalling.md)
		default:
			ai.IgnoredTypes[typedItem.ObjType] = lo.ValueOr(ai.IgnoredTypes, typedItem.ObjType, 0) + 1
		}
	}

	// Read the end of the array
	if _, err := dec.ReadToken(); err != nil {
		return fmt.Errorf("reading end of MdAllItems item object: %w", err)
	}

	return nil
}
