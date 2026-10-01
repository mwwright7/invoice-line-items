package lineitem

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLineItemMarshalJSON(t *testing.T) {
	li := LineItem{Quantity: 3, Description: `Blue "widget" | large`, UnitPrice: 1250}
	got, err := json.Marshal(li)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"quantity":3,"description":"Blue \"widget\" | large","unit_price":"12.50"}`
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestLineItemJSONRoundTrip(t *testing.T) {
	items := []LineItem{
		{Quantity: 1, Description: "Installation fee", UnitPrice: 4500},
		{Quantity: 2, Description: "Free sample", UnitPrice: 0},
	}
	data, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	var back []LineItem
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if len(back) != len(items) {
		t.Fatalf("got %d items, want %d", len(back), len(items))
	}
	for i := range items {
		if back[i] != items[i] {
			t.Errorf("item %d: got %+v, want %+v", i, back[i], items[i])
		}
	}
}

func TestLineItemUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    LineItem
		wantErr string
	}{
		{"string price", `{"quantity":2,"description":"A","unit_price":"12.50"}`, LineItem{2, "A", 1250}, ""},
		{"number price", `{"quantity":2,"description":"A","unit_price":12.5}`, LineItem{2, "A", 1250}, ""},
		{"whole number price", `{"quantity":1,"description":"A","unit_price":7}`, LineItem{1, "A", 700}, ""},
		{"zero quantity", `{"quantity":0,"description":"A","unit_price":"1.00"}`, LineItem{}, "quantity must be positive"},
		{"missing description", `{"quantity":1,"unit_price":"1.00"}`, LineItem{}, "description must not be empty"},
		{"negative price", `{"quantity":1,"description":"A","unit_price":"-1.00"}`, LineItem{}, "unit price must not be negative"},
		{"too many decimals", `{"quantity":1,"description":"A","unit_price":"1.005"}`, LineItem{}, "at most 2 decimal places"},
		{"bad price", `{"quantity":1,"description":"A","unit_price":"abc"}`, LineItem{}, "invalid amount"},
		{"wrong type", `{"quantity":"1","description":"A","unit_price":"1.00"}`, LineItem{}, "cannot unmarshal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got LineItem
			err := json.Unmarshal([]byte(tt.in), &got)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %v, want one containing %q", err, tt.wantErr)
				}
				if got != (LineItem{}) {
					t.Errorf("target modified on error: %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
