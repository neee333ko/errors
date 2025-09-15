package errors

import (
	"reflect"
	"testing"
)

func TestCreateString(t *testing.T) {
	type Struct struct {
		items  []string
		theMap map[string]Empty
		want   String
	}

	tests := []Struct{
		{
			items:  []string{"1", "2", "3"},
			theMap: nil,
			want: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
		},
		{
			items: nil,
			theMap: map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			},
			want: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
		},
	}

	for _, tt := range tests {
		var g String

		if tt.theMap == nil {
			g = NewString(tt.items...)

			if !reflect.DeepEqual(g, tt.want) {
				t.Errorf("NewString([]string) error: \nexpect:%v\ngot:%v\n", tt.want, g)
			}
		} else {
			g = MapKeysToString(tt.theMap)

			if !reflect.DeepEqual(g, tt.want) {
				t.Errorf("MapKeysToString(any) error: \nexpect:%v\ngot:%v\n", tt.want, g)
			}
		}
	}
}

func TestOperation(t *testing.T) {
	type Struct struct {
		op   string
		a    String
		b    String
		want String
		res  bool
	}

	tests := []Struct{
		{
			op: "intersection",
			a: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			b: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
			}),
			want: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
			}),
			res: false,
		},
		{
			op: "difference",
			a: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			b: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
			}),
			want: String(map[string]Empty{
				"3": Empty{},
			}),
			res: false,
		},
		{
			op: "union",
			a: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			b: String(map[string]Empty{
				"1": Empty{},
				"4": Empty{},
			}),
			want: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
				"4": Empty{},
			}),
			res: false,
		},
		{
			op: "issuperset",
			a: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			b: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
			}),
			want: nil,
			res:  true,
		},
		{
			op: "isequal",
			a: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			b: String(map[string]Empty{
				"1": Empty{},
				"2": Empty{},
				"3": Empty{},
			}),
			want: nil,
			res:  true,
		},
	}

	for _, tt := range tests {
		switch tt.op {
		case "intersection":
			s := tt.a.InterSection(tt.b)

			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("InterSection(String) error: \nexpect:%v\ngot:%v\n", tt.want, s)
			}
		case "difference":
			s := tt.a.Difference(tt.b)

			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("Difference(String) error: \nexpect:%v\ngot:%v\n", tt.want, s)
			}
		case "union":
			s := tt.a.Union(tt.b)

			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("Difference(String) error: \nexpect:%v\ngot:%v\n", tt.want, s)
			}
		case "issuperset":
			rr := tt.a.IsSuperSet(tt.b)

			if rr != tt.res {
				t.Errorf("IsSuperSet(String) error: \nexpect:%v\ngot:%v\n", tt.res, rr)
			}
		case "isequal":
			rr := tt.a.IsEqual(tt.b)

			if rr != tt.res {
				t.Errorf("IsEqual(String) error: \nexpect:%v\ngot:%v\n", tt.res, rr)
			}
		}
	}
}
