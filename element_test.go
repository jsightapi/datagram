package datagram

import "testing"

func TestElement_truncate(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{
			"Empty string",
			"",
			4,
			"",
		},
		{
			"Too short",
			"a",
			4,
			"a",
		},
		{
			"Too short 2",
			"ab",
			4,
			"ab",
		},
		{
			"Too short 3",
			"abc",
			4,
			"abc",
		},
		{
			"4 characters",
			"abcd",
			4,
			"...",
		},
		{
			"5 characters",
			"abcde",
			4,
			"...",
		},
		{
			"6 characters",
			"abcdef",
			4,
			"...",
		},
		{
			"7 characters",
			"abcdefg",
			4,
			"...",
		},
		{
			"8 characters",
			"abcdefgh",
			4,
			"a...",
		},
		{
			"8 characters",
			"abcdefghi",
			4,
			"ab...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := element{tt.name, tt.s, true}
			e.truncate(tt.n)
			if e.value != tt.want {
				t.Errorf("truncate() = %#v, want %#v", e.value, tt.want)
			}
		})
	}
}
