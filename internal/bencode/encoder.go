package bencode

import (
	"fmt"
	"io"
	"sort"
)

var (
	ErrInvalidType = fmt.Errorf("invalid type")
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Reset(w io.Writer) {
	e.w = w
}

func (e *Encoder) Encode(v interface{}) error {
	switch i := v.(type) {
	case string:
		return e.EncodeString(i)
	case int64:
		return e.EncodeInt(i)
	case []interface{}:
		return e.EncodeList(i)
	case map[string]interface{}:
		return e.EncodeDict(i)
	default:
		return ErrInvalidType
	}
	return nil
}

func (e *Encoder) EncodeString(v string) error {
	_, err := fmt.Fprintf(e.w, "%d:%s", len(v), v)
	return err
}

func (e *Encoder) EncodeInt(v int64) error {
	_, err := fmt.Fprintf(e.w, "i%de", v)
	return err
}

func (e *Encoder) EncodeList(v []interface{}) error {
	_, err := e.w.Write([]byte{'l'})
	if err != nil {
		return err
	}
	for _, item := range v {
		err = e.Encode(item)
		if err != nil {
			return err
		}
	}
	_, err = e.w.Write([]byte{'e'})
	return err
}

func (e *Encoder) EncodeDict(v map[string]interface{}) error {
	_, err := e.w.Write([]byte{'d'})
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		err = e.EncodeString(key)
		if err != nil {
			return err
		}
		err = e.Encode(v[key])
		if err != nil {
			return err
		}
	}
	_, err = e.w.Write([]byte{'e'})
	return err
}
