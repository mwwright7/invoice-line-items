// Package lineitem parses and validates invoice line items and formats
// them back into a readable table.
//
// The input format is one line item per line:
//
//	quantity | description | unit price
//
// Example:
//
//	3 | Blue widget | 12.50
//	1 | Installation fee | 45.00
//
// A description that needs to contain a literal '|' can escape it as
// '\|'. A literal backslash is written as '\\'.
//
// Money is kept as whole cents (Cents) rather than float64 so that sums
// and formatted output never drift from what was typed in.
package lineitem

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Cents is a monetary amount expressed in whole cents.
type Cents int64

// LineItem is one billable line on an invoice.
type LineItem struct {
	Quantity    int
	Description string
	UnitPrice   Cents
}

// Amount returns the line total: quantity times unit price.
func (li LineItem) Amount() Cents {
	return Cents(li.Quantity) * li.UnitPrice
}

// Validate reports whether a LineItem's fields are individually sensible,
// regardless of how the item was constructed.
func Validate(li LineItem) error {
	if li.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if strings.TrimSpace(li.Description) == "" {
		return errors.New("description must not be empty")
	}
	if li.UnitPrice < 0 {
		return errors.New("unit price must not be negative")
	}
	return nil
}

// ParseLine parses a single "qty | description | unit price" record and
// validates the result. A '|' inside a field must be escaped as '\|'.
func ParseLine(line string) (LineItem, error) {
	fields := splitEscaped(line)
	if len(fields) != 3 {
		return LineItem{}, fmt.Errorf("expected 3 fields separated by '|', got %d", len(fields))
	}

	qtyStr := strings.TrimSpace(fields[0])
	desc := strings.TrimSpace(fields[1])
	priceStr := strings.TrimSpace(fields[2])

	qty, err := strconv.Atoi(qtyStr)
	if err != nil {
		return LineItem{}, fmt.Errorf("invalid quantity %q: %w", qtyStr, err)
	}

	price, err := ParseCents(priceStr)
	if err != nil {
		return LineItem{}, fmt.Errorf("invalid unit price %q: %w", priceStr, err)
	}

	li := LineItem{Quantity: qty, Description: desc, UnitPrice: price}
	if err := Validate(li); err != nil {
		return LineItem{}, err
	}
	return li, nil
}

// ParseDocument parses one line item per line, skipping blank lines. It
// returns every line item that parsed successfully along with an error
// describing every line that did not (nil if all lines were valid).
func ParseDocument(text string) ([]LineItem, error) {
	var items []LineItem
	var problems []string

	for i, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		li, err := ParseLine(line)
		if err != nil {
			problems = append(problems, fmt.Sprintf("line %d: %v", i+1, err))
			continue
		}
		items = append(items, li)
	}

	if len(problems) > 0 {
		return items, errors.New(strings.Join(problems, "; "))
	}
	return items, nil
}

// splitEscaped splits line on '|', treating a backslash as an escape
// character: '\|' yields a literal '|' instead of a field boundary, and
// '\\' yields a literal '\'. A backslash before any other character is
// dropped and the character is kept as-is. A trailing unmatched
// backslash is kept literally.
func splitEscaped(line string) []string {
	var fields []string
	var cur strings.Builder

	escaped := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case escaped:
			cur.WriteByte(c)
			escaped = false
		case c == '\\':
			escaped = true
		case c == '|':
			fields = append(fields, cur.String())
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	if escaped {
		cur.WriteByte('\\')
	}
	fields = append(fields, cur.String())
	return fields
}

// EscapeField escapes '|' and '\' in s so it can be written back into a
// single field of the "qty | description | unit price" line format
// without being mistaken for a field separator.
func EscapeField(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `|`, `\|`)
	return s
}

// ParseCents parses a decimal string such as "12.50" or "12" into whole
// cents. It accepts zero, one, or two digits after the decimal point.
func ParseCents(s string) (Cents, error) {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	parts := strings.SplitN(s, ".", 2)
	if parts[0] == "" {
		return 0, errors.New("missing whole part")
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid whole part: %w", err)
	}

	var frac int64
	if len(parts) == 2 {
		fracStr := parts[1]
		if len(fracStr) > 2 {
			return 0, fmt.Errorf("at most 2 decimal places allowed, got %q", parts[1])
		}
		for len(fracStr) < 2 {
			fracStr += "0"
		}
		frac, err = strconv.ParseInt(fracStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid fractional part: %w", err)
		}
	}

	total := whole*100 + frac
	if neg {
		total = -total
	}
	return Cents(total), nil
}

// String formats cents as a fixed two-decimal amount, e.g. "12.50".
func (c Cents) String() string {
	neg := c < 0
	v := int64(c)
	if neg {
		v = -v
	}
	sign := ""
	if neg {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%02d", sign, v/100, v%100)
}
