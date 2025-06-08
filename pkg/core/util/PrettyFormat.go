package util

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"slices"
	"unsafe"
)

type bufferedState struct {
	fmt.State
	buf bytes.Buffer
}

func newBufferedState(w fmt.State) *bufferedState {
	return &bufferedState{
		State: w,
	}
}

func (bs *bufferedState) Write(b []byte) (n int, err error) {
	return bs.buf.Write(b)
}

func (bs *bufferedState) String() string {
	return bs.buf.String()
}

type PrettyMap[T comparable, U any] map[T]U

func (pm PrettyMap[T, U]) Format(w fmt.State, v rune) {
	PrettyFormat(w, v, pm)
}

func PrettyFormat(w fmt.State, v rune, i any) {
	r := reflect.ValueOf(i)
	for r.Kind() == reflect.Pointer {
		r = r.Elem()
	}
	if v == 'v' {
		switch r.Kind() {
		case reflect.Struct:
			prettyFormatStruct(w, v, r)
			return
		case reflect.Map:
			prettyFormatMap(w, v, r)
			return
		}
	}
	fmt.Fprintf(w, fmt.FormatString(w, v), i)
}

func SprettyFormat(w fmt.State, v rune, i any) string {
	bs := newBufferedState(w)
	PrettyFormat(bs, v, i)
	return bs.String()
}

func prettyFormatMap(w fmt.State, v rune, r reflect.Value) {
	if r.IsNil() {
		fmt.Fprintf(w, fmt.FormatString(w, v), r.Interface())
		return
	}
	keys := r.MapKeys()
	keyStrings := map[reflect.Value]string{}
	for _, k := range keys {
		keyStrings[k] = SprettyFormat(w, v, k.Interface())
	}
	slices.SortFunc(keys, func(a, b reflect.Value) int {
		switch aa := a.Interface().(type) {
		case int64:
			return cmp.Compare(aa, b.Interface().(int64))
		}
		return cmp.Compare(keyStrings[a], keyStrings[b])
	})
	if w.Flag('#') {
		fmt.Fprintf(w, "%s{", r.Type().String())
	} else {
		io.WriteString(w, "map[")
	}
	for i, key := range keys {
		if i > 0 {
			if w.Flag('#') {
				io.WriteString(w, ", ")
			} else {
				io.WriteString(w, " ")
			}
		}
		fmt.Fprintf(w, "%s:", keyStrings[key])
		PrettyFormat(w, v, r.MapIndex(key).Interface())
	}
	if w.Flag('#') {
		io.WriteString(w, "}")
	} else {
		io.WriteString(w, "]")
	}

}

func prettyFormatStruct(w fmt.State, v rune, r reflect.Value) {
	switch i := r.Interface().(type) {
	case big.Int:
		if w.Flag('#') {
			fmt.Fprintf(w, "big.Int(%s)", i.String())
		} else {
			io.WriteString(w, i.String())
		}
	case big.Rat:
		i.Set(&i)
		if w.Flag('#') {
			fmt.Fprintf(w, "big.Rat(%s)", i.RatString())
		} else {
			io.WriteString(w, i.RatString())
		}
	default:
		if w.Flag('#') {
			io.WriteString(w, r.Type().String())
		}
		io.WriteString(w, "{")
		for j := range r.Type().NumField() {
			if j > 0 {
				if w.Flag('#') {
					io.WriteString(w, ", ")
				} else {
					io.WriteString(w, " ")
				}
			}
			if w.Flag('#') {
				fmt.Fprintf(w, "%s:", r.Type().Field(j).Name)
			}
			rf := forceGetField(r, j)
			PrettyFormat(w, v, rf.Interface())
		}
		io.WriteString(w, "}")
	}
}

// gets a struct field that can be read even if it's unexported
func forceGetField(r reflect.Value, i int) reflect.Value {
	copy := reflect.New(r.Type()).Elem()
	copy.Set(r)
	rf := copy.Field(i)
	return reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
}
