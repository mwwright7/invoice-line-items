package lineitem

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// jsonItem is the wire form of a LineItem. Field names are snake_case to
// match the CSV header.
type jsonItem struct {
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
	UnitPrice   Cents  `json:"unit_price"`
}

// MarshalJSON writes the unit price as a decimal string ("12.50") rather
// than raw cents, so the output matches what ParseLine and WriteCSV use.
func (li LineItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonItem{li.Quantity, li.Description, li.UnitPrice})
}

// UnmarshalJSON decodes a line item and runs Validate on it, so a LineItem
// that came from JSON has passed the same checks as one from ParseLine.
func (li *LineItem) UnmarshalJSON(data []byte) error {
	var j jsonItem
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	decoded := LineItem{Quantity: j.Quantity, Description: j.Description, UnitPrice: j.UnitPrice}
	if err := Validate(decoded); err != nil {
		return err
	}
	*li = decoded
	return nil
}

// MarshalJSON encodes c as a JSON string such as "12.50". A string is used
// instead of a number so that consumers parsing into float64 can't
// reintroduce rounding error.
func (c Cents) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON accepts either a string ("12.50") or a plain JSON number
// (12.5) and applies the same rules as ParseCents.
func (c *Cents) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	s := string(data)
	if len(data) > 0 && data[0] == '"' {
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
	}
	v, err := ParseCents(s)
	if err != nil {
		return fmt.Errorf("invalid amount %q: %w", s, err)
	}
	*c = v
	return nil
}
