package api

import (
	"fmt"

	"adeynack.net/lapiasse/pkg/platform/mapex"
)

var registerTypeConversion = mapex.NewBiDirectionalLookupOrPanic(map[RegisterType]uint{
	RegisterTypeAsset:       1,
	RegisterTypeBank:        2,
	RegisterTypeCard:        3,
	RegisterTypeExpense:     4,
	RegisterTypeIncome:      5,
	RegisterTypeInstitution: 6,
	RegisterTypeInvestment:  7,
	RegisterTypeLiability:   8,
	RegisterTypeLoan:        9,
})

// Implement DB scanner and valuer interfaces for RegisterType
func (rt *RegisterType) Scan(value any) error {
	strValue, ok := value.(uint)
	if !ok {
		return fmt.Errorf("scanning RegisterType: expected uint, got %T", value)
	}

	apiValue, ok := registerTypeConversion.FromRight(strValue)
	if !ok {
		return fmt.Errorf("scanning RegisterType: unknown value %q", strValue)
	}

	*rt = apiValue

	return nil
}

func (rt RegisterType) Value() (any, error) {
	dbValue, ok := registerTypeConversion.FromLeft(rt)
	if !ok {
		return nil, fmt.Errorf("getting DB value for RegisterType: unknown value %q", rt)
	}

	return dbValue, nil
}
