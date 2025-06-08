package prob

import (
	"cmp"
	"fmt"
	"math/big"
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ExampleInt int

type ComparableStruct struct {
	V string
}

type ReversedKey string

func (k ReversedKey) StringKey() string {
	b := strings.Builder{}
	asBytes := []byte(k)
	for i := len(asBytes) - 1; i >= 0; i-- {
		b.WriteByte(asBytes[i])
	}
	return b.String()
}

func keysSortingTestCase[V cmp.Ordered](t *testing.T, values []V) {
	t.Run(fmt.Sprintf("correctly sorts %s keys", reflect.TypeFor[V]().String()), func(t *testing.T) {
		entries := []EntryT[V]{}
		for _, v := range values {
			entries = append(entries, EntryT[V]{
				Key: v, Value: big.NewRat(1, int64(len(values))),
			})
		}
		rand.Shuffle(len(entries), func(i, j int) {
			entries[i], entries[j] = entries[j], entries[i]
		})
		d, err := FromEntries(entries)
		require.NoError(t, err)

		actual := []V{}
		for _, k := range d.Keys() {
			e, _ := d.Lookup(k)
			actual = append(actual, e.Key)
		}

		assert.Equal(t, values, actual)
	})
}

func TestDist(t *testing.T) {
	t.Run("Format()", func(t *testing.T) {
		t.Run("prints a human-friendly string when using fmt", func(t *testing.T) {
			d, err := FromMap(MapT[string]{
				"bar": big.NewRat(1, 2),
				"foo": big.NewRat(1, 2),
			})
			require.NoError(t, err)

			require.Equal(t,
				"{map[bar:bar foo:foo] map[bar:1/2 foo:1/2] 2}",
				fmt.Sprintf("%v", d),
			)
			require.Equal(t,
				"[{map[bar:bar foo:foo] map[bar:1/2 foo:1/2] 2}]",
				fmt.Sprintf("%v", []Dist[string]{d}),
			)
			require.Equal(t,
				fmt.Sprintf(
					"prob.Dist[string]{%s}",
					strings.Join([]string{
						"vmap:map[prob.Key]string{\"bar\":\"bar\", \"foo\":\"foo\"}",
						"pmap:map[prob.Key]*big.Rat{\"bar\":big.Rat(1/2), \"foo\":big.Rat(1/2)}",
						"keyFuncKind:2",
					}, ", "),
				),
				fmt.Sprintf("%#v", d),
			)
		})
	})

	t.Run("Iter()", func(t *testing.T) {
		t.Run("iterates through the entries in key order", func(t *testing.T) {
			d, err := FromMap(MapT[ExampleInt]{
				1:  big.NewRat(1, 4),
				2:  big.NewRat(1, 4),
				10: big.NewRat(1, 4),
				11: big.NewRat(1, 4),
			})
			require.NoError(t, err)

			entries := []EntryT[ExampleInt]{}
			for v, p := range d.Iter() {
				entries = append(entries, EntryT[ExampleInt]{Key: v, Value: p})
			}

			assert.Equal(t,
				[]EntryT[ExampleInt]{
					{Key: 1, Value: big.NewRat(1, 4)},
					{Key: 2, Value: big.NewRat(1, 4)},
					{Key: 10, Value: big.NewRat(1, 4)},
					{Key: 11, Value: big.NewRat(1, 4)},
				},
				entries,
			)
		})
	})

	t.Run("key()", func(t *testing.T) {
		// this shouldn't normally be reachable
		t.Run("panics when the keyFuncKind is invalid", func(t *testing.T) {
			assert.PanicsWithError(t, "invalid keyFuncKind: 0", func() {
				Dist[string]{}.key("foo")
			})
		})
	})

	t.Run("Keys()", func(t *testing.T) {
		keysSortingTestCase(t, []int{-1, 0, 2, 10})
		keysSortingTestCase(t, []int8{-1, 0, 2, 10})
		keysSortingTestCase(t, []int16{-1, 0, 2, 10})
		keysSortingTestCase(t, []int32{-1, 0, 2, 10})
		keysSortingTestCase(t, []int64{-1, 0, 2, 10})
		keysSortingTestCase(t, []uint{0, 1, 2, 10})
		keysSortingTestCase(t, []uint8{0, 1, 2, 10})
		keysSortingTestCase(t, []uint16{0, 1, 2, 10})
		keysSortingTestCase(t, []uint32{0, 1, 2, 10})
		keysSortingTestCase(t, []uint64{0, 1, 2, 10})
		keysSortingTestCase(t, []uintptr{0, 1, 2, 10})
		keysSortingTestCase(t, []float32{-1, 0, 2, 10})
		keysSortingTestCase(t, []float64{-1, 0, 2, 10})
		keysSortingTestCase(t, []string{"10", "2", "a", "b"})
		keysSortingTestCase(t, []ExampleInt{-1, 0, 2, 10})
		keysSortingTestCase(t, []ReversedKey{"a", "ba", "cba", "abc"})
	})

	t.Run("Median()", func(t *testing.T) {
		t.Run("returns the median for a basic example", func(t *testing.T) {
			d, err := FromMap(MapT[string]{
				"foo": big.NewRat(1, 3),
				"bar": big.NewRat(1, 3),
				"baz": big.NewRat(1, 3),
			})
			require.NoError(t, err)

			require.Equal(t, "baz", d.Median(cmp.Compare))
		})

		t.Run("returns the greater value when there's a tie", func(t *testing.T) {
			d, err := FromMap(MapT[string]{
				"fob": big.NewRat(1, 4),
				"foo": big.NewRat(1, 4),
				"bar": big.NewRat(1, 4),
				"baz": big.NewRat(1, 4),
			})
			require.NoError(t, err)

			require.Equal(t, "fob", d.Median(cmp.Compare))
		})
	})

	t.Run("Percentile()", func(t *testing.T) {
		t.Run("returns the first value for the 0th percentile", func(t *testing.T) {
			d, err := FromMap(MapT[string]{
				"foo": big.NewRat(1, 3),
				"bar": big.NewRat(1, 3),
				"baz": big.NewRat(1, 3),
			})
			require.NoError(t, err)

			require.Equal(t, "bar", d.Percentile(0.0, cmp.Compare))
		})

		t.Run("returns the last value for the 100th percentile", func(t *testing.T) {
			d, err := FromMap(MapT[string]{
				"foo": big.NewRat(1, 3),
				"bar": big.NewRat(1, 3),
				"baz": big.NewRat(1, 3),
			})
			require.NoError(t, err)

			require.Equal(t, "foo", d.Percentile(1.0, cmp.Compare))
		})
	})

	t.Run("StringKey()", func(t *testing.T) {
		t.Run("returns identical values for semantically-equal dists", func(t *testing.T) {
			a, err := FromMap(MapT[*ComparableStruct]{
				{V: "foo"}: big.NewRat(1, 2),
				{V: "bar"}: big.NewRat(1, 2),
			})
			require.NoError(t, err)
			b, err := FromMap(MapT[*ComparableStruct]{
				{V: "foo"}: big.NewRat(1, 2),
				{V: "bar"}: big.NewRat(1, 2),
			})
			require.NoError(t, err)

			assert.Equal(t, a.StringKey(), b.StringKey())
		})
		t.Run("returns different values for semantically-distinct dists", func(t *testing.T) {
			a, err := FromMap(MapT[*ComparableStruct]{
				{V: "foo"}: big.NewRat(2, 3),
				{V: "bar"}: big.NewRat(1, 3),
			})
			require.NoError(t, err)
			b, err := FromMap(MapT[*ComparableStruct]{
				{V: "foo"}: big.NewRat(1, 3),
				{V: "bar"}: big.NewRat(2, 3),
			})
			require.NoError(t, err)

			assert.NotEqual(t, a.StringKey(), b.StringKey())
		})
	})

	t.Run("validate()", func(t *testing.T) {
		t.Run("returns an error when keyFuncKind is invalid", func(t *testing.T) {
			d, err := FromConst(1)
			require.NoError(t, err)

			d.keyFuncKind = keyFuncInvalidLast
			_, err = d.validate()
			require.ErrorContains(t, err, "invalid keyFuncKind:")
		})
	})
}

func TestFromEntries(t *testing.T) {
	t.Run("returns an error if the probabilities don't sum to 1", func(t *testing.T) {
		_, err := FromEntries([]EntryT[string]{
			{Key: "foo", Value: big.NewRat(1, 3)},
			{Key: "bar", Value: big.NewRat(1, 3)},
		})
		assert.ErrorContains(t, err, "sum of all probabilities must be 1")
	})

	t.Run("returns an error for an invalid value type", func(t *testing.T) {
		_, err := FromEntries([]EntryT[[]string]{
			{Key: []string{"foo"}, Value: big.NewRat(1, 2)},
			{Key: []string{"bar"}, Value: big.NewRat(1, 2)},
		})
		assert.ErrorContains(t, err, "[]string does not satisfy")
	})

	t.Run("returns an error when given duplicate keys", func(t *testing.T) {
		_, err := FromEntries([]EntryT[string]{
			{Key: "foo", Value: big.NewRat(1, 2)},
			{Key: "foo", Value: big.NewRat(1, 2)},
		})
		assert.ErrorContains(t, err, "multiple values with same key:")
	})
}
