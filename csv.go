package lineitem

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var csvHeader = []string{"quantity", "description", "unit_price"}

// WriteCSV writes items as CSV with a header row of quantity, description,
// unit_price. Unit price is written as a decimal string ("12.50") rather
// than raw cents so the file is readable in a spreadsheet.
func WriteCSV(w io.Writer, items []LineItem) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(csvHeader); err != nil {
		return err
	}
	for _, li := range items {
		record := []string{
			strconv.Itoa(li.Quantity),
			li.Description,
			li.UnitPrice.String(),
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

// ReadCSV reads line items from CSV with a header row of quantity,
// description, unit_price. Like ParseDocument, it returns every row that
// parsed successfully along with an error describing every row that did
// not (nil if all rows were valid), so a few bad rows don't discard an
// otherwise good import.
func ReadCSV(r io.Reader) ([]LineItem, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1

	header, err := cr.Read()
	if err == io.EOF {
		return nil, errors.New("empty CSV: missing header row")
	}
	if err != nil {
		return nil, fmt.Errorf("reading header: %w", err)
	}
	if !equalHeader(header, csvHeader) {
		return nil, fmt.Errorf("expected header %v, got %v", csvHeader, header)
	}

	var items []LineItem
	var problems []string
	row := 1 // the header occupies row 1

	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		row++
		if err != nil {
			problems = append(problems, fmt.Sprintf("row %d: %v", row, err))
			continue
		}
		if len(record) != 3 {
			problems = append(problems, fmt.Sprintf("row %d: expected 3 fields, got %d", row, len(record)))
			continue
		}

		qty, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			problems = append(problems, fmt.Sprintf("row %d: invalid quantity %q: %v", row, record[0], err))
			continue
		}
		price, err := ParseCents(strings.TrimSpace(record[2]))
		if err != nil {
			problems = append(problems, fmt.Sprintf("row %d: invalid unit price %q: %v", row, record[2], err))
			continue
		}

		li := LineItem{Quantity: qty, Description: strings.TrimSpace(record[1]), UnitPrice: price}
		if err := Validate(li); err != nil {
			problems = append(problems, fmt.Sprintf("row %d: %v", row, err))
			continue
		}
		items = append(items, li)
	}

	if len(problems) > 0 {
		return items, errors.New(strings.Join(problems, "; "))
	}
	return items, nil
}

func equalHeader(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range want {
		if strings.TrimSpace(got[i]) != want[i] {
			return false
		}
	}
	return true
}
