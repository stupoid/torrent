package bencode

import (
	"bytes"
	"testing"
)

func TestEncodeString(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{"spam", "4:spam", nil},
		{"", "0:", nil},
		{"hello", "5:hello", nil},
		{"foo", "3:foo", nil},
		{"0123456789", "10:0123456789", nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.EncodeString(test.input)
			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if buf.String() != test.expected {
				t.Errorf("expected %q, got %q for input %q", test.expected, buf.String(), test.input)
			}
		})
	}
}

func TestEncodeInt(t *testing.T) {
	tests := []struct {
		input       int64
		expected    string
		expectedErr error
	}{
		{123, "i123e", nil},
		{-456, "i-456e", nil},
		{0, "i0e", nil},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.EncodeInt(test.input)
			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if buf.String() != test.expected {
				t.Errorf("expected %q, got %q for input %q", test.expected, buf.String(), test.input)
			}
		})
	}
}

func TestEncodeList(t *testing.T) {
	tests := []struct {
		input       []interface{}
		expected    string
		expectedErr error
	}{
		{[]interface{}{"spam", "eggs"}, "l4:spam4:eggse", nil},
		{[]interface{}{int64(123), int64(456)}, "li123ei456ee", nil},
		{[]interface{}{}, "le", nil},
		{
			[]interface{}{
				map[string]interface{}{
					"key1": "value1",
					"key2": int64(42),
				},
				map[string]interface{}{
					"key3": "value3",
				},
			},
			"ld4:key16:value14:key2i42eed4:key36:value3ee",
			nil,
		},
		{
			[]interface{}{
				map[string]interface{}{
					"nested": map[string]interface{}{
						"innerKey": "innerValue",
					},
				},
			},
			"ld6:nestedd8:innerKey10:innerValueeee",
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.EncodeList(test.input)
			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if buf.String() != test.expected {
				t.Errorf("expected %q, got %q for input %q", test.expected, buf.String(), test.input)
			}
		})
	}
}

func TestEncodeDict(t *testing.T) {
	tests := []struct {
		input       map[string]interface{}
		expected    string
		expectedErr error
	}{
		{
			map[string]interface{}{
				"info": map[string]interface{}{
					"files": []interface{}{
						map[string]interface{}{
							"length": int64(12345),
							"path":   []interface{}{"filename"},
						},
					},
					"name":         "example.torrent",
					"piece length": int64(262144),
				},
			},
			"d4:infod5:filesld6:lengthi12345e4:pathl8:filenameeee4:name15:example.torrent12:piece lengthi262144eee",
			nil,
		},
		{
			map[string]interface{}{
				"info": map[string]interface{}{
					"length":       int64(12345),
					"name":         "example.torrent",
					"piece length": int64(262144),
				},
			},
			"d4:infod6:lengthi12345e4:name15:example.torrent12:piece lengthi262144eee",
			nil,
		},
		{
			map[string]interface{}{
				"b": "valueB",
				"a": "valueA",
				"c": "valueC",
			},
			"d1:a6:valueA1:b6:valueB1:c6:valueCe",
			nil,
		},
		{
			map[string]interface{}{
				"z": "last",
				"m": "middle",
				"a": "first",
			},
			"d1:a5:first1:m6:middle1:z4:laste",
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			err := encoder.EncodeDict(test.input)
			if err != test.expectedErr {
				t.Errorf("expected error %v, got %v for input %q", test.expectedErr, err, test.input)
			}
			if buf.String() != test.expected {
				t.Errorf("expected %q, got %q for input %q", test.expected, buf.String(), test.input)
			}
		})
	}
}
