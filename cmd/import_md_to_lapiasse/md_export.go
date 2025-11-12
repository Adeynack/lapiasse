package main

// MdExport represents the root of a Moneydance export JSON file.
type MdExport struct {
	Metadata MdMetadata `json:"metadata"`
	AllItems MdAllItems `json:"all_items"`
}
