package lineitem

import (
	"fmt"
	"strings"
)

// Format renders line items as an aligned table with a trailing total
// row. It is a pure function: the same items always produce the same
// string, which makes it easy to test with plain string comparisons.
func Format(items []LineItem) string {
	if len(items) == 0 {
		return "(no line items)\n"
	}

	qtyW := len("Qty")
	descW := len("Description")
	priceW := len("Unit Price")
	amountW := len("Amount")

	type row struct{ qty, desc, price, amount string }
	rows := make([]row, len(items))
	var total Cents

	for i, li := range items {
		amount := li.Amount()
		total += amount
		rows[i] = row{
			qty:    fmt.Sprintf("%d", li.Quantity),
			desc:   li.Description,
			price:  li.UnitPrice.String(),
			amount: amount.String(),
		}
		qtyW = maxWidth(qtyW, len(rows[i].qty))
		descW = maxWidth(descW, len(rows[i].desc))
		priceW = maxWidth(priceW, len(rows[i].price))
		amountW = maxWidth(amountW, len(rows[i].amount))
	}

	var b strings.Builder
	writeRow(&b, qtyW, descW, priceW, amountW, "Qty", "Description", "Unit Price", "Amount")
	writeSeparator(&b, qtyW, descW, priceW, amountW)
	for _, r := range rows {
		writeRow(&b, qtyW, descW, priceW, amountW, r.qty, r.desc, r.price, r.amount)
	}
	writeSeparator(&b, qtyW, descW, priceW, amountW)
	writeRow(&b, qtyW, descW, priceW, amountW, "", "Total", "", total.String())

	return b.String()
}

func writeRow(b *strings.Builder, qtyW, descW, priceW, amountW int, qty, desc, price, amount string) {
	fmt.Fprintf(b, "%*s  %-*s  %*s  %*s\n", qtyW, qty, descW, desc, priceW, price, amountW, amount)
}

func writeSeparator(b *strings.Builder, qtyW, descW, priceW, amountW int) {
	b.WriteString(strings.Repeat("-", qtyW+descW+priceW+amountW+6))
	b.WriteByte('\n')
}

func maxWidth(a, b int) int {
	if a > b {
		return a
	}
	return b
}
