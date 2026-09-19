package lineitem

import (
	"strings"
	"testing"
)

func TestWriteCSVThenReadCSVRoundTrips(t *testing.T) {
	items := []LineItem{
		{Quantity: 3, Description: "Blue widget", UnitPrice: 1250},
		{Quantity: 1, Description: "Installation fee", UnitPrice: 4500},
	}

	var b strings.Builder
	if err := WriteCSV(&b, items); err != nil {
		t.Fatalf("WriteCSV: unexpected error: %v", err)
	}

	got, err := ReadCSV(strings.NewReader(b.String()))
	if err != nil {
		t.Fatalf("ReadCSV: unexpected error: %v", err)
	}
	if len(got) != len(items) {
		t.Fatalf("got %d items, want %d", len(got), len(items))
	}
	for i := range items {
		if got[i] != items[i] {
			t.Errorf("item %d: got %+v, want %+v", i, got[i], items[i])
		}
	}
}

func TestWriteCSVEscapesDescriptionWithComma(t *testing.T) {
	items := []LineItem{{Quantity: 1, Description: "Widget, deluxe", UnitPrice: 500}}

	var b strings.Builder
	if err := WriteCSV(&b, items); err != nil {
		t.Fatalf("WriteCSV: unexpected error: %v", err)
	}

	got, err := ReadCSV(strings.NewReader(b.String()))
	if err != nil {
		t.Fatalf("ReadCSV: unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Description != "Widget, deluxe" {
		t.Fatalf("got %+v, want description %q", got, "Widget, deluxe")
	}
}

func TestReadCSVRejectsWrongHeader(t *testing.T) {
	_, err := ReadCSV(strings.NewReader("qty,desc,price\n1,Widget,5.00\n"))
	if err == nil {
		t.Fatal("expected error for wrong header")
	}
}

func TestReadCSVRejectsEmptyInput(t *testing.T) {
	_, err := ReadCSV(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestReadCSVCollectsRowErrors(t *testing.T) {
	doc := "quantity,description,unit_price\n" +
		"2,Widget,5.00\n" +
		"0,Bad quantity,5.00\n" +
		"1,Gadget,19.99\n"

	items, err := ReadCSV(strings.NewReader(doc))
	if err == nil {
		t.Fatal("expected error")
	}
	if len(items) != 2 {
		t.Fatalf("got %d valid items, want 2", len(items))
	}
	if !strings.Contains(err.Error(), "row 3") {
		t.Fatalf("error should mention row 3: %v", err)
	}
}
