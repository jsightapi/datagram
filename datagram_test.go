package datagram

import (
	"bytes"
	"net"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestData_truncatableElement(t *testing.T) {
	tests := []struct {
		name      string
		datagram  Datagram
		wantIndex int
	}{
		{
			"Empty set",
			New(),
			-1,
		},
		{
			"Only ellipsis",
			Datagram{
				list: []*element{
					{"str", ellipsis, true},
					{"str", ellipsis, true},
				},
			},
			1,
		},
		{
			"Untruncatable",
			Datagram{
				list: []*element{
					{"int", "1234", false},
					{"ip", "1.2.3.4", false},
				},
			},
			-1,
		},
		{
			"One truncatable element",
			Datagram{
				list: []*element{
					{"str", "abcd", true},
				},
			},
			0,
		},
		{
			"Two truncatable elements",
			Datagram{
				list: []*element{
					{"str", "abcd", true},
					{"str", "abcd", true},
				},
			},
			1,
		},
		{
			"Mixed",
			Datagram{
				list: []*element{
					{"str", "abcd", true},
					{"int", "1234", false},
					{"str", "abcd", true},
					{"int", "1234", false},
				},
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.datagram.truncatableElement()
			if e == nil && tt.wantIndex == -1 {
				return
			}
			if e == nil || tt.wantIndex == -1 || e != tt.datagram.list[tt.wantIndex] {
				t.Errorf("truncatableElement() = %v, wantIndex %v", e, tt.wantIndex)
			}
		})
	}
}

func Test_Pack(t *testing.T) { //nolint:funlen
	s100 := "1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 "

	tests := []struct {
		name     string
		datagram Datagram
		want     string
	}{
		{
			"Empty",
			New(),
			"\n\n",
		},
		{
			"One element",
			Datagram{
				list: []*element{
					{"aa", "abc", true},
				},
			},
			"aa\nabc\n",
		},
		{
			"Unicode",
			Datagram{
				list: []*element{
					{"aa", "Hello, 世界", true},
				},
			},
			"aa\n\"Hello, 世界\"\n",
		},
		{
			"Punctuation marks",
			Datagram{
				list: []*element{
					{"aa", "a,b;c.d-e", true},
				},
			},
			"aa\n\"a,b;c.d-e\"\n",
		},
		{
			"int64",
			Datagram{
				list: []*element{
					{"int64", strconv.FormatInt(-123, 10), false},
				},
			},
			"int64\n-123\n",
		},
		{
			"ip4",
			Datagram{
				list: []*element{
					{"ip", net.IPv4(1, 2, 3, 4).String(), false},
				},
			},
			"ip\n1.2.3.4\n",
		},
		{
			"Non-unique keys",
			Datagram{
				list: []*element{
					{"aa", "abc", false},
					{"bb", "-123", false},
					{"bb", "1.2.3.4", false},
				},
			},
			"aa,bb,bb\nabc,-123,1.2.3.4\n",
		},
		{
			"Too long string",
			Datagram{
				list: []*element{
					{"aa", strings.Repeat(s100, 10), true},
				},
			},
			"aa\n1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 123456...\n", //nolint:lll
		},
		{
			"Couple of too long strings",
			Datagram{
				list: []*element{
					{"aa", strings.Repeat(s100, 10), true},
					{"bb", strings.Repeat(s100, 10), true},
				},
			},
			"aa,bb\n1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890...,...\n", //nolint:lll
		},
		{
			"Truncatable & untruncatable",
			Datagram{
				list: []*element{
					{"i1", "1", false},
					{"aa", strings.Repeat(s100, 10), true},
					{"i2", "2", false},
					{"bb", strings.Repeat(s100, 10), true},
					{"i3", "3", false},
				},
			},
			"i1,aa,i2,bb,i3\n1,1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 1234567890 123456...,2,...,3\n", //nolint:lll
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.datagram.Pack()
			if len(got) > maxSize {
				t.Errorf("Encode() too long string, len = %d", len(got))
			}
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("Encode() got = %q,\n\twant %q", got, tt.want)
			}
		})
	}
}

func TestUnpack(t *testing.T) { //nolint:funlen
	type args struct {
		s string
	}
	t.Run("positive", func(t *testing.T) {
		tests := []struct {
			name string
			args args
			want Datagram
		}{
			{
				"One key-value",
				args{
					"key\nvalue",
				},
				Datagram{
					list: []*element{
						{"key", "value", false},
					},
				},
			},
			{
				"Two key-values",
				args{
					"key1,key2\nvalue1,value2",
				},
				Datagram{
					list: []*element{
						{"key1", "value1", false},
						{"key2", "value2", false},
					},
				},
			},
			{
				"Extra new line character",
				args{
					"key1,key2\n\nvalue1,value2",
				},
				Datagram{
					list: []*element{
						{"key1", "value1", false},
						{"key2", "value2", false},
					},
				},
			},
			{
				"Missing one of values",
				args{
					"key1,key2\nvalue1",
				},
				Datagram{
					list: []*element{
						{"key1", "value1", false},
						{"key2", "", false},
					},
				},
			},
			{
				"Missing one of keys",
				args{
					"key1\nvalue1,value2",
				},
				Datagram{
					list: []*element{
						{"key1", "value1", false},
					},
				},
			},
		}
		wantErr := false
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Unpack([]byte(tt.args.s))
				if (err != nil) != wantErr {
					t.Errorf("Unpack() error = %v, wantErr %v", err, wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Unpack() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		tests := []struct {
			name string
			args args
		}{
			{
				"Empty CSV",
				args{
					"",
				},
			},
			{
				"Missing all values",
				args{
					"key1,key2\n",
				},
			},
			{
				"Missing all values 2",
				args{
					"key1,key2",
				},
			},
		}
		wantErr := true
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := Unpack([]byte(tt.args.s))
				if (err != nil) != wantErr {
					t.Errorf("Unpack() error = %v, wantErr %v", err, wantErr)
					return
				}
			})
		}
	})
}
