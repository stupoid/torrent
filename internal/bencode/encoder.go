package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

type encoder struct {
	bytes.Buffer
}

func (e *encoder) writeByteString(s string) {
	e.WriteString(fmt.Sprintf("%d:%s", len(s), s))
}

func (e *encoder) writeInteger(i int64) {
	e.WriteString(fmt.Sprintf("i%de", i))
}
func (e *encoder) writeUnsignedInteger(i uint64) {
	e.WriteString(fmt.Sprintf("i%de", i))
}

func (e *encoder) writeInterfaceType(i interface{}) {
	switch v := i.(type) {
	case string:
		e.writeByteString(v)
	case int, int8, int16, int32, int64:
		e.writeInteger(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		e.writeUnsignedInteger(reflect.ValueOf(v).Uint())
	case []interface{}:
		e.writeList(v)
	case map[string]interface{}:
		e.writeDictionary(v)
	}
}

func (e *encoder) writeList(l []interface{}) {
	e.WriteByte('l')
	for _, v := range l {
		e.writeInterfaceType(v)
	}
	e.WriteByte('e')
}

func (e *encoder) writeDictionary(d map[string]interface{}) {
	keys := make(sort.StringSlice, len(d))
	i := 0
	for k := range d {
		keys[i] = k
		i++
	}
	keys.Sort()

	e.WriteByte('d')
	for _, k := range keys {
		e.writeByteString(k)
		e.writeInterfaceType(d[k])
	}
	e.WriteByte('e')
}

func Encode(v interface{}) []byte {
	e := &encoder{}
	e.writeInterfaceType(v)
	return e.Bytes()
}
