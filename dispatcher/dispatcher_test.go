package dispatcher

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"mysql> select 2;\b\b\b\b\bt 2;\u001B[K\b\b\b\bct 2;", "mysql> select 2;"},
		{"mysql> select 2;\b\b\b\b\b\b\u001B[Pc\b\u001B[1@ec", "mysql> select 2;"},
		{"mysql> show databases;\u001B[9Gelect 1;\u001B[K", "mysql> select 1;"},
		{"mysql> show tables;\u001B[9Gelect 1;\u001B[K", "mysql> select 1;"},
		{"mysql> select 1;\u001B[9Ghow databases;", "mysql> show databases;"},
		{"mysql> show databases;\u001B[10Gaaa;\u001B[K\b\b\b\bow databases;\u001B[9Geaaa;\u001B[K\u001B[9Ghow databases;\u001B[9Gelect 1;\u001B[K", "mysql> select 1;"},
		{"mysql> show databases;\u001B[10Gaaa;\u001B[K\b\b\b\bow databases;\u001B[9Geaaa;\u001B[K\u001B[9Ghow databases;\u001B[9Gelect 1;\u001B[K", "mysql> select 1;"},
		{"mysql> a;\b\bdatabases;\u001B[8Ga;\u001B[K\b\bdatabases;\u001B[8Ga;\u001B[K\b\bdatabases;\u001B[8Ga;\u001B[K", "mysql> a;"},
	}
	for _, test := range tests {
		got := Parse(strings.NewReader(test.in))
		if got != test.want {
			t.Errorf("Parse(%q) == %q, want %q", test.in, got, test.want)
		}
	}
}