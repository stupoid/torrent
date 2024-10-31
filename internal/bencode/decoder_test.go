package bencode

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"4:spam", "spam", nil},
		{"0:", "", nil},
		{"5:hello", "hello", nil},
		{"3:foo", "foo", nil},
		{"10:0123456789", "0123456789", nil},
		{"4:sp", "", ErrReadValueFailed},   // Error case: length mismatch
		{"4spam", "", ErrReadLengthFailed}, // Error case: missing colon
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			reader := bufio.NewReader(strings.NewReader(test.input))
			decoder := NewDecoder(reader)
			result, err := decoder.DecodeString()

			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if result != test.expected {
				t.Errorf("expected %q, got %q for input %q", test.expected, result, test.input)
			}
		})
	}
}
func TestDecodeInt(t *testing.T) {
	tests := []struct {
		input       string
		expected    int64
		expectedErr error
	}{
		{"i123e", 123, nil},
		{"i0e", 0, nil},
		{"i-456e", -456, nil},
		{"i123", 0, ErrInvalidEndingByte},  // Error case: missing ending 'e'
		{"123e", 0, ErrInvalidLeadingByte}, // Error case: missing leading 'i'
		{"i12a3e", 0, ErrReadValueFailed},  // Error case: invalid integer format
		{"", 0, ErrReadLeadingFailed},      // Error case: empty input
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			reader := bufio.NewReader(strings.NewReader(test.input))
			decoder := NewDecoder(reader)
			result, err := decoder.DecodeInt()

			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if result != test.expected {
				t.Errorf("expected %d, got %d for input %q", test.expected, result, test.input)
			}
		})
	}
}
func TestDecodeList(t *testing.T) {
	tests := []struct {
		input       string
		expected    []interface{}
		expectedErr error
	}{
		{"li123ei456ee", []interface{}{int64(123), int64(456)}, nil},
		{"l4:spam4:eggse", []interface{}{"spam", "eggs"}, nil},
		{"le", []interface{}{}, nil},
		{"li123e4:spame", []interface{}{int64(123), "spam"}, nil},
		{"ld3:foo3:baree", []interface{}{map[string]interface{}{"foo": "bar"}}, nil},
		{"li123ei456e", nil, ErrInvalidEndingByte},  // Error case: missing ending 'e'
		{"i123ei456ee", nil, ErrInvalidLeadingByte}, // Error case: missing leading 'l'
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			reader := bufio.NewReader(strings.NewReader(test.input))
			decoder := NewDecoder(reader)
			result, err := decoder.DecodeList()

			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v for input %q", test.expected, result, test.input)
			}
		})
	}
}

func TestDecodeDict(t *testing.T) {
	tests := []struct {
		input       string
		expected    map[string]interface{}
		expectedErr error
	}{
		{"d3:bar4:spam3:fooi42ee", map[string]interface{}{"bar": "spam", "foo": int64(42)}, nil},
		{"d3:food3:bar3:bazee", map[string]interface{}{"foo": map[string]interface{}{"bar": "baz"}}, nil},
		{"d3:fool3:bar3:bazee", map[string]interface{}{"foo": []interface{}{"bar", "baz"}}, nil},
		{"de", map[string]interface{}{}, nil},
		{"d3:bar4:spame", map[string]interface{}{"bar": "spam"}, nil},
		{"d3:bar4:spam3:fooi42e", nil, ErrInvalidEndingByte},  // Error case: missing ending 'e'
		{"3:bar4:spam3:fooi42ee", nil, ErrInvalidLeadingByte}, // Error case: missing leading 'd'
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			reader := bufio.NewReader(strings.NewReader(test.input))
			decoder := NewDecoder(reader)
			result, err := decoder.DecodeDict()

			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected %v, got %v for input %q", test.expected, result, test.input)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		input       string
		expected    interface{}
		expectedErr error
	}{
		{
			"d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:\"Debian CD from cdimage.debian.org\"13:creation datei1391870037e9:httpseedsl85:http://cdimage.debian.org/cdimage/release/7.4.0/iso-cd/debian-7.4.0-amd64-netinst.iso85:http://cdimage.debian.org/cdimage/archive/7.4.0/iso-cd/debian-7.4.0-amd64-netinst.isoe4:infod6:lengthi232783872e4:name30:debian-7.4.0-amd64-netinst.iso12:piece lengthi262144e6:pieces0:ee",
			map[string]interface{}{
				"announce":      "http://bttracker.debian.org:6969/announce",
				"comment":       "\"Debian CD from cdimage.debian.org\"",
				"creation date": int64(1391870037),
				"httpseeds": []interface{}{
					"http://cdimage.debian.org/cdimage/release/7.4.0/iso-cd/debian-7.4.0-amd64-netinst.iso",
					"http://cdimage.debian.org/cdimage/archive/7.4.0/iso-cd/debian-7.4.0-amd64-netinst.iso",
				},
				"info": map[string]interface{}{
					"length":       int64(232783872),
					"name":         "debian-7.4.0-amd64-netinst.iso",
					"piece length": int64(262144),
					"pieces":       "",
				},
			},
			nil,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			reader := bufio.NewReader(strings.NewReader(test.input))
			decoder := NewDecoder(reader)
			result, err := decoder.Decode()
			if err != test.expectedErr {
				t.Fatalf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Fatalf("expected %v, got %v for input %q", test.expected, result, test.input)
			}
		})
	}
}
