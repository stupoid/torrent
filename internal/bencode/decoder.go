package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// Custom error types
var (
	ErrInvalidEndingByte   = errors.New("invalid ending byte")
	ErrInvalidLeadingByte  = errors.New("invalid leading byte")
	ErrInvalidLengthFormat = errors.New("invalid length format")
	ErrReadLeadingFailed   = errors.New("failed to leading byte")
	ErrReadLengthFailed    = errors.New("failed to read length")
	ErrReadValueFailed     = errors.New("failed to read value")
)

func Decode(data []byte) (interface{}, error) {
	return nil, nil
}

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r *bufio.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) Reset(r *bufio.Reader) {
	d.r = r
}

func (d Decoder) Decode() (interface{}, error) {
	leading, err := d.r.Peek(1)
	if err != nil {
		return nil, ErrInvalidLeadingByte
	}
	switch leading[0] {
	case 'i':
		return d.DecodeInt()
	case 'l':
		return d.DecodeList()
	case 'd':
		return d.DecodeDict()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.DecodeString()
	default:
		return nil, ErrInvalidLeadingByte
	}
}

func (d Decoder) DecodeString() (string, error) {
	lenStr, err := d.r.ReadString(':')
	if err != nil {
		return "", ErrReadLengthFailed
	}
	length, err := strconv.Atoi(lenStr[:len(lenStr)-1])
	if err != nil {
		return "", ErrInvalidLengthFormat
	}
	value := make([]byte, length)
	if _, err := io.ReadFull(d.r, value); err != nil {
		return "", ErrReadValueFailed
	}
	return string(value), nil
}

func (d Decoder) DecodeInt() (int64, error) {
	leading, err := d.r.ReadByte()
	if err != nil {
		return 0, ErrReadLeadingFailed
	}
	if leading != 'i' {
		return 0, ErrInvalidLeadingByte
	}
	valueString, err := d.r.ReadString('e')
	if err != nil {
		if err == io.EOF {
			return 0, ErrInvalidEndingByte
		}
		return 0, ErrReadValueFailed
	}
	value, err := strconv.ParseInt(valueString[:len(valueString)-1], 10, 64)
	if err != nil {
		return 0, ErrReadValueFailed
	}
	return value, nil
}

func (d Decoder) DecodeList() ([]interface{}, error) {
	leading, err := d.r.ReadByte()
	if err != nil {
		return nil, ErrReadLeadingFailed
	}
	if leading != 'l' {
		return nil, ErrInvalidLeadingByte
	}
	list := make([]interface{}, 0)
	for {
		next, err := d.r.Peek(1)
		if err != nil {
			return nil, ErrInvalidEndingByte
		}
		if next[0] == 'e' {
			d.r.ReadByte()
			break
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		list = append(list, value)
	}
	return list, nil
}

func (d Decoder) DecodeDict() (map[string]interface{}, error) {
	leading, err := d.r.ReadByte()
	if err != nil {
		return nil, ErrReadLeadingFailed
	}
	if leading != 'd' {
		return nil, ErrInvalidLeadingByte
	}
	dict := make(map[string]interface{})
	for {
		next, err := d.r.Peek(1)
		if err != nil {
			return nil, ErrInvalidEndingByte
		}
		if next[0] == 'e' {
			d.r.ReadByte()
			break
		}
		key, err := d.DecodeString()
		if err != nil {
			return nil, err
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		dict[key] = value
	}
	return dict, nil
}
