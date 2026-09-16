package lineitem

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	li, err := ParseLine("3 | Blue widget | 12.50")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := LineItem{Quantity: 3, Description: "Blue widget", UnitPrice: 1250}
	if li != want {
		t.Fatalf("got %+v, want %+v", li, want)
	}
	if got := li.Amount(); got != 3750 {
		t.Fatalf("Amount() = %d, want 3750", got)
	}
}

func TestParseLineRejectsBadInput(t *testing.T) {
	cases := []string{
		"0 | Widget | 1.00",  // non-positive quantity
		"1 |  | 1.00",        // empty description
		"1 | Widget | -1.00", // negative price
		"1 | Widget | 1.005", // too many decimal places
		"1 | Widget",         // wrong field count
	}
	for _, c := range cases {
		if _, err := ParseLine(c); err == nil {
			t.Errorf("ParseLine(%q): expected error, got none", c)
		}
	}
}

func TestParseLineEscapedPipeInDescription(t *testing.T) {
	li, err := ParseLine(`2 | Widget \| Deluxe | 5.00`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "Widget | Deluxe"
	if li.Description != want {
		t.Fatalf("Description = %q, want %q", li.Description, want)
	}
}

func TestParseLineEscapedBackslashInDescription(t *testing.T) {
	li, err := ParseLine(`1 | C:\Program Files\\ | 5.00`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `C:Program Files\`
	if li.Description != want {
		t.Fatalf("Description = %q, want %q", li.Description, want)
	}
}

func TestEscapeFieldRoundTrips(t *testing.T) {
	desc := `Widget | Deluxe \ Edition`
	line := fmt.Sprintf("1 | %s | 5.00", EscapeField(desc))
	li, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if li.Description != desc {
		t.Fatalf("Description = %q, want %q", li.Description, desc)
	}
}

func TestParseDocument(t *testing.T) {
	doc := "2 | Widget | 5.00\n\n1 | Gadget | 19.99\n"
	items, err := ParseDocument(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
}

func TestParseDocumentCollectsErrors(t *testing.T) {
	doc := "2 | Widget | 5.00\nnot a line item\n1 | Gadget | 19.99\n"
	items, err := ParseDocument(doc)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(items) != 2 {
		t.Fatalf("got %d valid items, want 2", len(items))
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Fatalf("error should mention line 2: %v", err)
	}
}

func TestFormatIncludesTotal(t *testing.T) {
	items := []LineItem{
		{Quantity: 2, Description: "Widget", UnitPrice: 500},
		{Quantity: 1, Description: "Gadget", UnitPrice: 1999},
	}
	out := Format(items)
	if !strings.Contains(out, "Total") {
		t.Fatalf("expected Total row in output:\n%s", out)
	}
	if !strings.Contains(out, "29.99") {
		t.Fatalf("expected total amount 29.99 in output:\n%s", out)
	}
}
