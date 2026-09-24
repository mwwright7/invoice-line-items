package lineitem

import (
	"strings"
	"testing"
)

func TestFormatCurrencyPrefixesMoneyColumns(t *testing.T) {
	items := []LineItem{
		{Quantity: 2, Description: "Widget", UnitPrice: 500},
		{Quantity: 1, Description: "Gadget", UnitPrice: 1999},
	}
	out := FormatCurrency(items, "$")

	if !strings.Contains(out, "$5.00") {
		t.Fatalf("expected unit price $5.00 in output:\n%s", out)
	}
	if !strings.Contains(out, "$10.00") {
		t.Fatalf("expected amount $10.00 in output:\n%s", out)
	}
	if !strings.Contains(out, "$29.99") {
		t.Fatalf("expected total $29.99 in output:\n%s", out)
	}
	if strings.Contains(out, "Qty") == false {
		t.Fatalf("expected header row in output:\n%s", out)
	}
}

func TestFormatCurrencyEmptySymbolMatchesFormat(t *testing.T) {
	items := []LineItem{
		{Quantity: 3, Description: "Blue widget", UnitPrice: 1250},
	}
	if got, want := FormatCurrency(items, ""), Format(items); got != want {
		t.Fatalf("FormatCurrency with empty symbol = %q, want %q", got, want)
	}
}

func TestFormatCurrencyNoItems(t *testing.T) {
	if got, want := FormatCurrency(nil, "$"), "(no line items)\n"; got != want {
		t.Fatalf("FormatCurrency(nil) = %q, want %q", got, want)
	}
}
