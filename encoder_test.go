package bone

import (
	"bytes"
	"testing"
)

func FuzzEncode(f *testing.F) {
	// Basic block values
	f.Add([]byte{0x20})                                                 // Simple B0 value (false)
	f.Add([]byte{0x21})                                                 // Simple B0 value (true)
	f.Add([]byte{0x10})                                                 // Simple B0 value (int 0)
	f.Add([]byte{0x11})                                                 // Simple B0 value (int 1)
	f.Add([]byte{0xFF, 0x21})                                           // Level extension
	f.Add([]byte{0xFF, 0xFF, 0x21})                                     // Double level extension
	f.Add([]byte{0x0F, 0xAA})                                           // B1 value (negative int)
	f.Add([]byte{0x0E, 0x11, 0x22})                                     // B2 value (negative int)
	f.Add([]byte{0x0D, 0x11, 0x22, 0x33})                               // B3 value (negative int)
	f.Add([]byte{0x0C, 0x11, 0x22, 0x33, 0x44})                         // B4 value (negative int)
	f.Add([]byte{0x18, 0xAA})                                           // B1 value (positive int)
	f.Add([]byte{0x19, 0x11, 0x22})                                     // B2 value (positive int)
	f.Add([]byte{0x30, 0xAA})                                           // B1 value (type 0x30)
	f.Add([]byte{0x40, 0xAA, 0xBB})                                     // B2 value (type 0x40)
	f.Add([]byte{0x50, 0xAA, 0xBB, 0xCC})                               // B3 value (type 0x50)
	f.Add([]byte{0x70, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // B8 value (float 1.0)

	// String values
	f.Add([]byte{0x90, 0x00})                         // Empty string
	f.Add([]byte{0x91, 0x41, 0x42, 0x43, 0x00})       // String "ABC"
	f.Add([]byte{0x92, 0x41, 0x00, 0x01, 0x43, 0x00}) // String with escaped null byte
	f.Add([]byte{0x90, 0x41, 0x42, 0x43, 0x00})       // String with different code
	f.Add([]byte{0x91, 0x00, 0x01, 0x00, 0x01, 0x00}) // String with multiple escaped nulls

	// Tuple values
	f.Add([]byte{0xA0, 0x20})                   // T1 (boolean)
	f.Add([]byte{0xB0, 0x20, 0x21})             // T2 (two booleans)
	f.Add([]byte{0xC0, 0x20, 0x21, 0x10})       // T3 (various values)
	f.Add([]byte{0xD0, 0x20, 0x21, 0x10, 0x11}) // T4 (various values)

	// List values
	f.Add([]byte{0xF0, 0x00})                                                 // Empty list
	f.Add([]byte{0xF1, 0x00})                                                 // Empty list (different code)
	f.Add([]byte{0xF0, 0x10, 0x20, 0x30, 0xAA, 0x00})                         // List with values
	f.Add([]byte{0xF0, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x00}) // List with many B0 values

	// Complex nested values
	f.Add([]byte{
		0xF0,                         // List
		0x91, 0x41, 0x42, 0x43, 0x00, // String "ABC"
		0xA1, 0x20, // T1 with boolean
		0xB1, 0x0F, 0xAA, 0xFF, 0x21, // T2 with values
		0x00, // List terminator
	})

	// Deeply nested structures
	f.Add([]byte{
		0xF0,       // Outer list
		0xF1,       // Inner list 1
		0xA0, 0x21, // T1 (true)
		0xB0, 0x20, 0x21, // T2 (false, true)
		0x00,                         // End inner list 1
		0xF2,                         // Inner list 2
		0x91, 0x58, 0x59, 0x5A, 0x00, // String "XYZ"
		0x0F, 0xFF, // B1 value
		0x00, // End inner list 2
		0x00, // End outer list
	})

	// Mixed type list with all major value types
	f.Add([]byte{
		0xF0,       // List
		0x20,       // B0 (false)
		0x21,       // B0 (true)
		0x10,       // B0 (int 0)
		0x0F, 0x7F, // B1 (negative int)
		0x18, 0x7F, // B1 (positive int)
		0x30, 0xCC, // B1 (type 0x30)
		0x70, 0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18, // B8 (float 3.14)
		0x90, 0x74, 0x65, 0x73, 0x74, 0x00, // String "test"
		0xA0, 0x21, // T1 (true)
		0xB0, 0x20, 0x21, // T2 (false, true)
		0xF1, 0x10, 0x11, 0x00, // Nested list
		0x00, // End of list
	})

	// Multiple extension levels in various positions
	f.Add([]byte{
		0xFF, 0xFF, 0xFF, 0x20, // Triple level extension B0 (false)
	})

	f.Add([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0x21, // Quadruple level extension B0 (true)
	})

	f.Add([]byte{
		0xFF, 0xFF, 0x0F, 0xCC, // Double level extension B1 value
	})

	f.Add([]byte{
		0xFF, 0xFF, 0xFF, 0x91, 0x41, 0x42, 0x43, 0x00, // Triple level extension String "ABC"
	})

	// Complex structure with multiple extension levels
	f.Add([]byte{
		0xF0,       // Outer list
		0xFF, 0x21, // Level 1 extension (true)
		0xFF, 0xFF, 0x20, // Level 2 extension (false)
		0xFF, 0xFF, 0xFF, 0x21, // Level 3 extension (true)
		0xA0, 0xFF, 0x21, // T1 with level 1 extension
		0xB0, 0xFF, 0xFF, 0x20, 0xFF, 0x21, // T2 with level 2 and level 1 extensions
		0x00, // End list
	})

	// Nested list with different extension levels
	f.Add([]byte{
		0xF0,       // Outer list
		0xF1,       // Inner list 1
		0xFF, 0x20, // Level 1 extension (false)
		0xFF, 0xFF, 0x21, // Level 2 extension (true)
		0x00,                                     // End inner list 1
		0xF2,                                     // Inner list 2
		0xFF, 0xFF, 0xFF, 0x90, 0x41, 0x42, 0x00, // Level 3 extension string "AB"
		0xFF, 0x0F, 0xAA, // Level 1 extension B1 value
		0x00, // End inner list 2
		0x00, // End outer list
	})

	// Alternating extension levels with block values of different sizes
	f.Add([]byte{
		0xF0,             // List
		0xFF, 0x0F, 0xAA, // Level 1 extension B1
		0xFF, 0xFF, 0x0E, 0xBB, 0xCC, // Level 2 extension B2
		0xFF, 0xFF, 0xFF, 0x0D, 0xDD, 0xEE, 0xFF, // Level 3 extension B3
		0xFF, 0x30, 0x55, // Level 1 extension B1 (type 0x30)
		0xFF, 0xFF, 0x40, 0x66, 0x77, // Level 2 extension B2 (type 0x40)
		0x00, // End list
	})

	// Mixed extension levels in tuple structures
	f.Add([]byte{
		0xA0, 0xFF, 0xFF, 0x21, // T1 with level 2 extension
	})

	f.Add([]byte{
		0xB0, 0xFF, 0x20, 0xFF, 0xFF, 0xFF, 0x21, // T2 with level 1 and level 3 extensions
	})

	f.Add([]byte{
		0xC0,       // T3
		0xFF, 0x20, // Level 1 extension (false)
		0xFF, 0xFF, 0x21, // Level 2 extension (true)
		0xFF, 0xFF, 0xFF, 0x10, // Level 3 extension (int 0)
	})
	f.Fuzz(func(t *testing.T, payload []byte) {
		values, err := Decode(payload)
		if err != nil {
			return
		}
		if !bytes.Equal(Encode(values), payload) {
			t.Errorf("re-encoded bytes do not match the original payload")
		}
	})
}
