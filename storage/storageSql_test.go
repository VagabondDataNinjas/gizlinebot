package storage

import (
	"testing"
)

func TestSanitizePrice(t *testing.T) {
	var data = []struct {
		input       string
		expected    float64
		expectedErr bool
	}{
		// valid inputs
		{"34.56", 34.56, false},
		{"  12  ", 12.0, false},
		{"  000.01  ", 0.01, false},
		{"sdfgd34dg", 34.0, false},
		{"ลิตรละ33บาทค่ะ", 33.0, false},
		{"ราคาน้ำมันลิตรละ 27.95 บาท", 27.95, false},
		// multiple dot chars outside the digits
		{"ราคาน้ำมันลิต...รละ 27.95.,.. บ.าท", 27.95, false},
		// invalid inputs,
		{"", 0.0, true},
		{" ", 0.0, true},
		{"ราคาน้ำมันลิตรละ", 0.0, true},
		{"27.9.5", 0.0, true},
	}

	for _, item := range data {
		price, err := sanitizePrice(item.input)
		if err != nil && !item.expectedErr {
			t.Errorf("Got unexpected error %s", err)
		}
		if price != item.expected {
			t.Errorf("Unexpected \"%.2f\", expected \"%.2f\"", price, item.expected)
		}
	}
}
