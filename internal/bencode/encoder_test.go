package bencode

import (
	"math"
	"testing"
)

func TestEncoder(t *testing.T) {
	t.Run("writeByteString", func(t *testing.T) {
		testCases := map[string]struct {
			input    string
			expected string
		}{
			"normal": {
				input:    "foo",
				expected: "3:foo",
			},
			"empty": {
				input:    "",
				expected: "0:",
			},
			"special characters": {
				input:    "foo:bar",
				expected: "7:foo:bar",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeByteString(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})

	t.Run("writeInteger", func(t *testing.T) {
		testCases := map[string]struct {
			input    int64
			expected string
		}{
			"positive": {
				input:    42,
				expected: "i42e",
			},
			"negative": {
				input:    -42,
				expected: "i-42e",
			},
			"zero": {
				input:    0,
				expected: "i0e",
			},
			"max positive": {
				input:    math.MaxInt64,
				expected: "i9223372036854775807e",
			},
			"max negative": {
				input:    math.MinInt64,
				expected: "i-9223372036854775808e",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeInteger(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})

	t.Run("writeUnsignedInteger", func(t *testing.T) {
		testCases := map[string]struct {
			input    uint64
			expected string
		}{
			"small": {
				input:    42,
				expected: "i42e",
			},
			"zero": {
				input:    0,
				expected: "i0e",
			},
			"max": {
				input:    math.MaxUint64,
				expected: "i18446744073709551615e",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeUnsignedInteger(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})

	t.Run("writeInterface handle all integer types", func(t *testing.T) {
		testCases := map[string]struct {
			input    interface{}
			expected string
		}{
			"int": {
				input:    int(42),
				expected: "i42e",
			},
			"int8": {
				input:    int8(42),
				expected: "i42e",
			},
			"int16": {
				input:    int16(42),
				expected: "i42e",
			},
			"int32": {
				input:    int32(42),
				expected: "i42e",
			},
			"int64": {
				input:    int64(42),
				expected: "i42e",
			},
			"uint": {
				input:    uint(42),
				expected: "i42e",
			},
			"uint8": {
				input:    uint8(42),
				expected: "i42e",
			},
			"uint16": {
				input:    uint16(42),
				expected: "i42e",
			},
			"uint32": {
				input:    uint32(42),
				expected: "i42e",
			},
			"uint64": {
				input:    uint64(42),
				expected: "i42e",
			},
			"max int": {
				input:    int64(math.MaxInt64),
				expected: "i9223372036854775807e",
			},
			"max negative int": {
				input:    int64(math.MinInt64),
				expected: "i-9223372036854775808e",
			},
			"negative int": {
				input:    int(-42),
				expected: "i-42e",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeInterfaceType(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})

	t.Run("writeList", func(t *testing.T) {
		testCases := map[string]struct {
			input    []interface{}
			expected string
		}{
			"normal": {
				input:    []interface{}{"foo", 42},
				expected: "l3:fooi42ee",
			},
			"empty": {
				input:    []interface{}{},
				expected: "le",
			},
			"nested empty list": {
				input:    []interface{}{[]interface{}{}},
				expected: "llee",
			},
			"mixed types": {
				input:    []interface{}{"foo", 42, []interface{}{"bar", 84}},
				expected: "l3:fooi42el3:bari84eee",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeList(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})

	t.Run("writeDictionary", func(t *testing.T) {
		testCases := map[string]struct {
			input    map[string]interface{}
			expected string
		}{
			"normal": {
				input: map[string]interface{}{
					"foo": "bar",
					"baz": 42,
				},
				expected: "d3:bazi42e3:foo3:bare",
			},
			"nested": {
				input: map[string]interface{}{
					"foo": "bar",
					"baz": map[string]interface{}{
						"qux":  42,
						"fizz": []interface{}{"buzz", "buzzz"},
					},
					"eee": []interface{}{},
				},
				expected: "d3:bazd4:fizzl4:buzz5:buzzze3:quxi42ee3:eeele3:foo3:bare",
			},
			"empty": {
				input:    map[string]interface{}{},
				expected: "de",
			},
			"nested empty dictionary": {
				input: map[string]interface{}{
					"foo": map[string]interface{}{},
				},
				expected: "d3:foodee",
			},
			"empty string value": {
				input: map[string]interface{}{
					"foo": "",
				},
				expected: "d3:foo0:e",
			},
			"keys sorted as raw strings": {
				input: map[string]interface{}{
					"1":  "one",
					"10": "ten",
					"2":  "two",
				},
				expected: "d1:13:one2:103:ten1:23:twoe",
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				e := &encoder{}
				e.writeDictionary(tc.input)

				got := e.String()
				if got != tc.expected {
					t.Fatalf("expected %q, got %q", tc.expected, got)
				}
			})
		}
	})
}
