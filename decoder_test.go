package bone

import (
	"fmt"
	"testing"
)

func TestDecodeIllegal(t *testing.T) {
	t.Run("Base illegal type codes", func(t *testing.T) {
		for i := range 8 {
			tc := byte(i)
			t.Run(fmt.Sprintf("typecode 0x%02X", tc), func(t *testing.T) {
				decoder := &Decoder{}
				err := decoder.Accept(tc)
				if err == nil {
					t.Fatalf("Expected error for illegal type code 0x%02X, got nil", tc)
				}
			})
		}
	})

	t.Run("Illegal extension sequences", func(t *testing.T) {
		illegalExtensions := []struct {
			name  string
			bytes []byte
		}{
			{
				name:  "Level extension followed by illegal type code (0x00)",
				bytes: []byte{0xFF, 0x00},
			},
			{
				name:  "Level extension followed by illegal type code (0x07)",
				bytes: []byte{0xFF, 0x07},
			},
			{
				name:  "Level extension followed by illegal extension type code (0x0F)",
				bytes: []byte{0xFF, 0x0F},
			},
			{
				name:  "Multiple level extension followed by illegal extension type code (0x1F)",
				bytes: []byte{0xFF, 0xFF, 0xFF, 0x1F},
			},
			{
				name:  "Multiple level extensions followed by illegal type code",
				bytes: []byte{0xFF, 0xFF, 0x00},
			},
		}
		for _, tc := range illegalExtensions {
			t.Run(tc.name, func(t *testing.T) {
				decoder := &Decoder{}
				var err error
				for i := range len(tc.bytes) - 1 {
					err = decoder.Accept(tc.bytes[i])
					if err != nil {
						t.Fatalf("Unexpected error during setup: %v", err)
					}
				}
				lastByte := tc.bytes[len(tc.bytes)-1]
				err = decoder.Accept(lastByte)
				if err == nil {
					t.Fatalf("Expected error for illegal extension sequence ending with 0x%02X, got nil", lastByte)
				}
			})
		}
	})
}

func TestDecodeBlock(t *testing.T) {
	payload := []byte{
		// B0
		0x10, 0x17,
		0x20,
		0xFF, 0x21,
		0xFF, 0xFF, 0x2F,
		// B1
		0x0F, 0xAA,
		0x18, 0xBB,
		0x30, 0xCC,
		0xFF, 0x31, 0xDD,
		0xFF, 0xFF, 0x3F, 0xEE,
		// B2
		0x0E, 0x11, 0x22,
		0x19, 0x33, 0x44,
		0x40, 0x55, 0x66,
		0xFF, 0x41, 0x77, 0x88,
		0xFF, 0xFF, 0x4F, 0x99, 0xAA,
		// B3
		0x0D, 0x11, 0x22, 0x33,
		0x1A, 0x44, 0x55, 0x66,
		0x50, 0x77, 0x88, 0x99,
		0xFF, 0x51, 0xAA, 0xBB, 0xCC,
		0xFF, 0xFF, 0x5F, 0xDD, 0xEE, 0xFF,
		// B4
		0x0C, 0x11, 0x22, 0x33, 0x44,
		0x1B, 0x55, 0x66, 0x77, 0x88,
		0x60, 0x99, 0xAA, 0xBB, 0xCC,
		0xFF, 0x61, 0xDD, 0xEE, 0xFF, 0x00,
		0xFF, 0xFF, 0x6F, 0x11, 0x22, 0x33, 0x44,
		// B5
		0x0B, 0x01, 0x02, 0x03, 0x04, 0x05,
		0x1C, 0x06, 0x07, 0x08, 0x09, 0x0A,
		// B6
		0x0A, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x1D, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
		// B7
		0x09, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x1E, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
		// B8
		0x08, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x1F, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x70, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0xFF, 0x71, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
		// B16
		0x80, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
		0xFF, 0x81, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
	}
	expected := []*Value{
		// B0
		{Code: 0x10, Level: 0},
		{Code: 0x17, Level: 0},
		{Code: 0x20, Level: 0},
		{Code: 0x21, Level: 1},
		{Code: 0x2F, Level: 2},
		// B1
		{Code: 0x0F, Level: 0, Bytes: []byte{0xAA}},
		{Code: 0x18, Level: 0, Bytes: []byte{0xBB}},
		{Code: 0x30, Level: 0, Bytes: []byte{0xCC}},
		{Code: 0x31, Level: 1, Bytes: []byte{0xDD}},
		{Code: 0x3F, Level: 2, Bytes: []byte{0xEE}},
		// B2
		{Code: 0x0E, Level: 0, Bytes: []byte{0x11, 0x22}},
		{Code: 0x19, Level: 0, Bytes: []byte{0x33, 0x44}},
		{Code: 0x40, Level: 0, Bytes: []byte{0x55, 0x66}},
		{Code: 0x41, Level: 1, Bytes: []byte{0x77, 0x88}},
		{Code: 0x4F, Level: 2, Bytes: []byte{0x99, 0xAA}},
		// B3
		{Code: 0x0D, Level: 0, Bytes: []byte{0x11, 0x22, 0x33}},
		{Code: 0x1A, Level: 0, Bytes: []byte{0x44, 0x55, 0x66}},
		{Code: 0x50, Level: 0, Bytes: []byte{0x77, 0x88, 0x99}},
		{Code: 0x51, Level: 1, Bytes: []byte{0xAA, 0xBB, 0xCC}},
		{Code: 0x5F, Level: 2, Bytes: []byte{0xDD, 0xEE, 0xFF}},
		// B4
		{Code: 0x0C, Level: 0, Bytes: []byte{0x11, 0x22, 0x33, 0x44}},
		{Code: 0x1B, Level: 0, Bytes: []byte{0x55, 0x66, 0x77, 0x88}},
		{Code: 0x60, Level: 0, Bytes: []byte{0x99, 0xAA, 0xBB, 0xCC}},
		{Code: 0x61, Level: 1, Bytes: []byte{0xDD, 0xEE, 0xFF, 0x00}},
		{Code: 0x6F, Level: 2, Bytes: []byte{0x11, 0x22, 0x33, 0x44}},
		// B5
		{Code: 0x0B, Level: 0, Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05}},
		{Code: 0x1C, Level: 0, Bytes: []byte{0x06, 0x07, 0x08, 0x09, 0x0A}},
		// B6
		{Code: 0x0A, Level: 0, Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}},
		{Code: 0x1D, Level: 0, Bytes: []byte{0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C}},
		// B7
		{Code: 0x09, Level: 0, Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
		{Code: 0x1E, Level: 0, Bytes: []byte{0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E}},
		// B8
		{Code: 0x08, Level: 0, Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
		{Code: 0x1F, Level: 0, Bytes: []byte{0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}},
		{Code: 0x70, Level: 0, Bytes: []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}},
		{Code: 0x71, Level: 1, Bytes: []byte{0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}},
		// B16
		{Code: 0x80, Level: 0, Bytes: []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30}},
		{Code: 0x81, Level: 1, Bytes: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40}},
	}

	values, err := Decode(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	for i, exp := range expected {
		if i >= len(values) {
			t.Fatalf("Missing value at index %d", i)
		}
		val := values[i]
		if val.Code != exp.Code {
			t.Errorf("Value %d: Expected type code 0x%02X, got 0x%02X", i, exp.Code, val.Code)
		}
		if val.Level != exp.Level {
			t.Errorf("Value %d: Expected level %d, got %d", i, exp.Level, val.Level)
		}
		if exp.Bytes != nil {
			if len(val.Bytes) != len(exp.Bytes) {
				t.Errorf("Value %d: Expected %d bytes, got %d bytes", i, len(exp.Bytes), len(val.Bytes))
			} else {
				for j, b := range exp.Bytes {
					if val.Bytes[j] != b {
						t.Errorf("Value %d: Byte at index %d: expected 0x%02X, got 0x%02X", i, j, b, val.Bytes[j])
					}
				}
			}
		} else if len(val.Bytes) > 0 {
			t.Errorf("Value %d: Expected no bytes, got %d bytes", i, len(val.Bytes))
		}
	}
}

func TestDecodeString(t *testing.T) {
	payload := []byte{
		0x90, 0x00,
		0x9F, 0x00,
		0xFF, 0x90, 0x00,
		0xFF, 0xFF, 0x9F, 0x00,
		0x91, 0x00, 0x01, 0x00,
		0x92, 0xAA, 0xBB, 0x00, 0x01, 0x91, 0x01, 0xFF, 0x00,
	}
	expected := []*Value{
		{Code: 0x90, Level: 0},
		{Code: 0x9F, Level: 0},
		{Code: 0x90, Level: 1},
		{Code: 0x9F, Level: 2},
		{Code: 0x91, Level: 0, Bytes: []byte{0x00}},
		{Code: 0x92, Level: 0, Bytes: []byte{0xAA, 0xBB, 0x00, 0x91, 0x01, 0xFF}},
	}

	values, err := Decode(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	for i, exp := range expected {
		if i >= len(values) {
			t.Fatalf("Missing value at index %d", i)
		}
		val := values[i]
		if val.Code != exp.Code {
			t.Errorf("Value %d: Expected type code 0x%02X, got 0x%02X", i, exp.Code, val.Code)
		}
		if val.Level != exp.Level {
			t.Errorf("Value %d: Expected level %d, got %d", i, exp.Level, val.Level)
		}
		if len(val.Bytes) != len(exp.Bytes) {
			t.Errorf("Value %d: Expected %d bytes, got %d bytes", i, len(exp.Bytes), len(val.Bytes))
		} else {
			for j, b := range exp.Bytes {
				if val.Bytes[j] != b {
					t.Errorf("Value %d: Byte at index %d: expected 0x%02X, got 0x%02X", i, j, b, val.Bytes[j])
				}
			}
		}
	}
}

func TestDecodeTuple(t *testing.T) {
	payload := []byte{
		// T1
		0xA0, 0x20,
		0xFF, 0xAF, 0xFF, 0x21,
		// T2
		0xB0, 0x20, 0x21,
		0xFF, 0xBF, 0x10, 0xFF, 0xFF, 0x2F,
		// T3
		0xC0, 0x20, 0x30, 0xAA, 0x40, 0xBB, 0xCC,
		0xFF, 0xCF, 0x20, 0x30, 0xDD, 0xFF, 0x40, 0xEE, 0xFF,
		// T4
		0xD0, 0x20, 0x30, 0xAA, 0x40, 0xBB, 0xCC, 0x50, 0xDD, 0xEE, 0xFF,
		0xFF, 0xDF, 0x21, 0x31, 0x11, 0x41, 0x22, 0x33, 0x51, 0x44, 0x55, 0x66,
	}

	expected := []*Value{
		// T1
		{
			Code:  0xA0,
			Level: 0,
			Values: []*Value{
				{Code: 0x20, Level: 0},
			},
		},
		{
			Code:  0xAF,
			Level: 1,
			Values: []*Value{
				{Code: 0x21, Level: 1},
			},
		},
		// T2
		{
			Code:  0xB0,
			Level: 0,
			Values: []*Value{
				{Code: 0x20, Level: 0},
				{Code: 0x21, Level: 0},
			},
		},
		{
			Code:  0xBF,
			Level: 1,
			Values: []*Value{
				{Code: 0x10, Level: 0},
				{Code: 0x2F, Level: 2},
			},
		},
		// T3
		{
			Code:  0xC0,
			Level: 0,
			Values: []*Value{
				{Code: 0x20, Level: 0},
				{Code: 0x30, Level: 0, Bytes: []byte{0xAA}},
				{Code: 0x40, Level: 0, Bytes: []byte{0xBB, 0xCC}},
			},
		},
		{
			Code:  0xCF,
			Level: 1,
			Values: []*Value{
				{Code: 0x20, Level: 0},
				{Code: 0x30, Level: 0, Bytes: []byte{0xDD}},
				{Code: 0x40, Level: 1, Bytes: []byte{0xEE, 0xFF}},
			},
		},
		// T4
		{
			Code:  0xD0,
			Level: 0,
			Values: []*Value{
				{Code: 0x20, Level: 0},
				{Code: 0x30, Level: 0, Bytes: []byte{0xAA}},
				{Code: 0x40, Level: 0, Bytes: []byte{0xBB, 0xCC}},
				{Code: 0x50, Level: 0, Bytes: []byte{0xDD, 0xEE, 0xFF}},
			},
		},
		{
			Code:  0xDF,
			Level: 1,
			Values: []*Value{
				{Code: 0x21, Level: 0},
				{Code: 0x31, Level: 0, Bytes: []byte{0x11}},
				{Code: 0x41, Level: 0, Bytes: []byte{0x22, 0x33}},
				{Code: 0x51, Level: 0, Bytes: []byte{0x44, 0x55, 0x66}},
			},
		},
	}

	values, err := Decode(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	for i, exp := range expected {
		if i >= len(values) {
			t.Fatalf("Missing value at index %d", i)
		}
		val := values[i]

		if val.Code != exp.Code {
			t.Errorf("Value %d: Expected type code 0x%02X, got 0x%02X", i, exp.Code, val.Code)
		}

		if val.Level != exp.Level {
			t.Errorf("Value %d: Expected level %d, got %d", i, exp.Level, val.Level)
		}

		if len(val.Values) != len(exp.Values) {
			t.Errorf("Value %d: Expected %d nested values, got %d", i, len(exp.Values), len(val.Values))
			continue
		}

		for j, expNested := range exp.Values {
			if j >= len(val.Values) {
				t.Errorf("Value %d: Missing nested value at index %d", i, j)
				continue
			}

			nested := val.Values[j]

			if nested.Code != expNested.Code {
				t.Errorf("Value %d, nested %d: Expected type code 0x%02X, got 0x%02X", i, j, expNested.Code, nested.Code)
			}

			if nested.Level != expNested.Level {
				t.Errorf("Value %d, nested %d: Expected level %d, got %d", i, j, expNested.Level, nested.Level)
			}

			if len(nested.Bytes) != len(expNested.Bytes) {
				t.Errorf("Value %d, nested %d: Expected %d bytes, got %d bytes", i, j, len(expNested.Bytes), len(nested.Bytes))
			} else if expNested.Bytes != nil {
				for k, b := range expNested.Bytes {
					if nested.Bytes[k] != b {
						t.Errorf("Value %d, nested %d: Byte at index %d: expected 0x%02X, got 0x%02X", i, j, k, b, nested.Bytes[k])
					}
				}
			}
		}
	}
}

func TestDecodeList(t *testing.T) {

	payload := []byte{
		0xF0, 0x00, // L0
		0xFF, 0xFE, 0xFF, 0x20, 0x00, // L1
		0xF0, 0x10, 0x20, 0x30, 0xAA, 0x40, 0xBB, 0xCC, 0x00, // L4
	}

	expected := []*Value{
		// L0 - Empty list
		{
			Code:   0xF0,
			Level:  0,
			Values: []*Value{},
		},
		// L1 with level extensions
		{
			Code:  0xFE,
			Level: 1,
			Values: []*Value{
				{Code: 0x20, Level: 1},
			},
		},
		// L4 - List with multiple values
		{
			Code:  0xF0,
			Level: 0,
			Values: []*Value{
				{Code: 0x10, Level: 0},
				{Code: 0x20, Level: 0},
				{Code: 0x30, Level: 0, Bytes: []byte{0xAA}},
				{Code: 0x40, Level: 0, Bytes: []byte{0xBB, 0xCC}},
			},
		},
	}

	values, err := Decode(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	for i, exp := range expected {
		if i >= len(values) {
			t.Fatalf("Missing value at index %d", i)
		}
		val := values[i]

		if val.Code != exp.Code {
			t.Errorf("Value %d: Expected type code 0x%02X, got 0x%02X", i, exp.Code, val.Code)
		}

		if val.Level != exp.Level {
			t.Errorf("Value %d: Expected level %d, got %d", i, exp.Level, val.Level)
		}

		if len(val.Values) != len(exp.Values) {
			t.Errorf("Value %d: Expected %d list values, got %d", i, len(exp.Values), len(val.Values))
			continue
		}

		for j, expNested := range exp.Values {
			if j >= len(val.Values) {
				t.Errorf("Value %d: Missing list value at index %d", i, j)
				continue
			}

			nested := val.Values[j]

			if nested.Code != expNested.Code {
				t.Errorf("Value %d, nested %d: Expected type code 0x%02X, got 0x%02X", i, j, expNested.Code, nested.Code)
			}

			if nested.Level != expNested.Level {
				t.Errorf("Value %d, nested %d: Expected level %d, got %d", i, j, expNested.Level, nested.Level)
			}

			// Check bytes in nested value
			if expNested.Bytes != nil {
				if len(nested.Bytes) != len(expNested.Bytes) {
					t.Errorf("Value %d, nested %d: Expected %d bytes, got %d bytes", i, j, len(expNested.Bytes), len(nested.Bytes))
				} else {
					for k, b := range expNested.Bytes {
						if nested.Bytes[k] != b {
							t.Errorf("Value %d, nested %d: Byte at index %d: expected 0x%02X, got 0x%02X", i, j, k, b, nested.Bytes[k])
						}
					}
				}
			} else if len(nested.Bytes) > 0 {
				t.Errorf("Value %d, nested %d: Expected no bytes, got %d bytes", i, j, len(nested.Bytes))
			}
		}
	}
}

func TestDecodeNestedMix(t *testing.T) {
	payload := []byte{
		0xF0,
		0x91, 0x00, 0x01, 0x00,
		0xA1,
		0xB1,
		0xF1, 0x10, 0x00,
		0xFF, 0xAF, 0xFF, 0x9F, 0xFF, 0x00,
		0xFF, 0xFF, 0xFE, 0x00,
		0x00,
	}

	expected := []*Value{
		{
			Code:  0xF0,
			Level: 0,
			Values: []*Value{
				{Code: 0x91, Level: 0, Bytes: []byte{0x00}},
				{
					Code:  0xA1,
					Level: 0,
					Values: []*Value{
						{Code: 0xB1, Level: 0, Values: []*Value{
							{
								Code:  0xF1,
								Level: 0,
								Values: []*Value{
									{Code: 0x10, Level: 0},
								},
							},
							{
								Code:  0xAF,
								Level: 1,
								Values: []*Value{
									{Code: 0x9F, Level: 1, Bytes: []byte{0xFF}},
								},
							},
						}},
					},
				},
				{Code: 0xFE, Level: 2, Values: []*Value{}},
			},
		},
	}

	values, err := Decode(payload)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	if len(values) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(values))
	}

	// Helper function to compare nested values recursively
	var compareValues func(t *testing.T, val, exp *Value, path string)
	compareValues = func(t *testing.T, val, exp *Value, path string) {
		if val.Code != exp.Code {
			t.Errorf("%s: Expected type code 0x%02X, got 0x%02X", path, exp.Code, val.Code)
		}

		if val.Level != exp.Level {
			t.Errorf("%s: Expected level %d, got %d", path, exp.Level, val.Level)
		}

		// Compare bytes
		if len(val.Bytes) != len(exp.Bytes) {
			t.Errorf("%s: Expected %d bytes, got %d bytes", path, len(exp.Bytes), len(val.Bytes))
		} else if exp.Bytes != nil {
			for k, b := range exp.Bytes {
				if val.Bytes[k] != b {
					t.Errorf("%s: Byte at index %d: expected 0x%02X, got 0x%02X", path, k, b, val.Bytes[k])
				}
			}
		}

		// Compare nested values
		if len(val.Values) != len(exp.Values) {
			t.Errorf("%s: Expected %d nested values, got %d", path, len(exp.Values), len(val.Values))
			return
		}

		for j, expNested := range exp.Values {
			if j >= len(val.Values) {
				t.Errorf("%s: Missing nested value at index %d", path, j)
				continue
			}

			nested := val.Values[j]
			nestedPath := fmt.Sprintf("%s.Values[%d]", path, j)
			compareValues(t, nested, expNested, nestedPath)
		}
	}

	for i, exp := range expected {
		if i >= len(values) {
			t.Fatalf("Missing value at index %d", i)
		}
		val := values[i]
		path := fmt.Sprintf("values[%d]", i)
		compareValues(t, val, exp, path)
	}
}
