package csvquery

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
			[][]string{{"id", "name"}, {"1", "Alice"}, {"2", "Bob"}},
			"SELECT * FROM test WHERE id = '1'",
			[][]string{{"id", "name"}, {"1", "Alice"}},
		},
		{
			"function_test",
			[][]string{{"id", "name"}, {"1", "Alice"}, {"2", "Bob"}},
			"SELECT UPPER(name) AS name FROM function_test",
			[][]string{{"name"}, {"ALICE"}, {"BOB"}},
		},
		{
			"dashed-table-name",
			[][]string{{"id", "name"}, {"1", "Alice"}, {"2", "Bob"}},
			"SELECT id FROM 'dashed-table-name'",
			[][]string{{"id"}, {"1"}, {"2"}},
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
