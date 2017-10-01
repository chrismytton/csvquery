package csvsql

import (
	"reflect"
	"testing"
)

func TestCSVTableCreateStatement(t *testing.T) {
	cases := []struct {
		tableName string
		records   [][]string
		query     string
		result    [][]string
	}{
		{
			"test",
			[][]string{[]string{"id", "name"}, []string{"1", "Alice"}, []string{"2", "Bob"}},
			"SELECT * FROM test WHERE id = '1'",
			[][]string{[]string{"id", "name"}, []string{"1", "Alice"}},
		},
	}
	for _, c := range cases {
		q, err := New()
		if err != nil {
			t.Error(err)
		}
		defer q.Close()
		err = q.Insert(c.tableName, c.records)
		if err != nil {
			t.Error(err)
		}
		got, err := q.Query(c.query)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, c.result) {
			t.Errorf("Query(%q) == %q, want %q", c.query, got, c.result)
		}
	}
}
