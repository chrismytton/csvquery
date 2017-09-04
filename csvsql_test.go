package csvsql

import (
	"testing"
)

func TestCSVTableCreateStatement(t *testing.T) {
	cases := []struct {
		tableName       string
		rows            [][]string
		createStatement string
	}{
		{
			"test-headers-only",
			[][]string{[]string{"id", "name"}},
			"CREATE TABLE test-headers-only (id, name)",
		},
		{
			"test-names",
			[][]string{[]string{"name"}, []string{"Alice"}, []string{"Bob"}},
			"CREATE TABLE test-names (name)",
		},
	}
	for _, c := range cases {
		table := &CSVTable{name: c.tableName, rows: c.rows}
		got := table.CreateStatement()
		if got != c.createStatement {
			t.Errorf("CreateStatement(%q, %q) == %q, want %q", c.tableName, c.rows, got, c.createStatement)
		}
	}
}
