package moneydance

// Export represents the root of a Moneydance export JSON file.
type Export struct {
	Metadata Metadata `json:"metadata"`
	AllItems AllItems `json:"all_items"`
}
