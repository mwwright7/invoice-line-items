# invoice-line-items

A small Go library for parsing invoice line items from plain text,
validating them, and printing them back out as an aligned table.

The problem this solves: invoice data usually starts life as something
loose and human-typed (a text field, a copy-pasted list, a quick note),
and it needs to become something you can trust before it goes into a
total. This package draws a hard line between "text someone typed" and
"a LineItem that has already been checked," and keeps money as integer
cents throughout so totals never drift from floating point rounding.

## Format

One line item per line:

```
quantity | description | unit price
```

Example:

```
3 | Blue widget | 12.50
1 | Installation fee | 45.00
```

Quantity must be a positive integer, description must be non-empty
after trimming whitespace, and unit price must be zero or positive
with at most two decimal places.

## Usage

```go
package main

import (
	"fmt"

	"github.com/mwwright7/invoice-line-items"
)

func main() {
	doc := `2 | Blue widget | 12.50
1 | Installation fee | 45.00`

	items, err := lineitem.ParseDocument(doc)
	if err != nil {
		fmt.Println("some lines failed to parse:", err)
	}

	fmt.Print(lineitem.Format(items))
}
```

Output:

```
Qty  Description        Unit Price  Amount
-------------------------------------------
  2  Blue widget             12.50   25.00
  1  Installation fee        45.00   45.00
-------------------------------------------
     Total                            70.00
```

## Design notes

- `ParseLine` and `ParseDocument` never mutate anything and never touch
  a filesystem or clock, so every case in the test suite is a plain
  input/output comparison.
- `Validate` is exported separately from `ParseLine` so callers building
  `LineItem` values by hand (not from text) can still run the same
  checks.
- Money is represented as `Cents` (an `int64`) rather than `float64` to
  avoid rounding surprises when items are summed.

## Status

Early skeleton. The line format is deliberately minimal right now; see
below for what's planned.

## Roadmap

- Escape or quote the `|` separator so descriptions can contain it
- CSV import/export
- Optional currency symbol in `Format` output
- JSON encoding for `LineItem`

## License

MIT, see [LICENSE](LICENSE).
