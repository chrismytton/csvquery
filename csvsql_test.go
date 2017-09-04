package csvsql

import "testing"

func TestCSVTableCreateStatement(t *testing.T) {
	cases := []struct {
		tableName       string
		records         [][]string
		createStatement string
		insertStatement string
	}{
		{
			"test-headers-only",
			[][]string{[]string{"id", "name"}},
			"CREATE TABLE test-headers-only (id, name)",
			"INSERT INTO test-headers-only (id, name) VALUES (?, ?)",
		},
		{
			"test-names",
			[][]string{[]string{"name"}, []string{"Alice"}, []string{"Bob"}},
			"CREATE TABLE test-names (name)",
			"INSERT INTO test-names (name) VALUES (?)",
		},
	}
	for _, c := range cases {
		table := &CSVTable{c.tableName, c.records}
		got := table.CreateStatement()
		if got != c.createStatement {
			t.Errorf("CreateStatement(%q, %q) == %q, want %q", c.tableName, c.records, got, c.createStatement)
		}
		got = table.InsertStatement()
		if got != c.insertStatement {
			t.Errorf("InsertStatement(%q, %q) == %q, want %q", c.tableName, c.records, got, c.insertStatement)
		}
	}
}
