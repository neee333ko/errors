package errors

import "testing"

func TestNewAggregate(t *testing.T) {
	type Struct struct {
		errs []error
		want string
	}

	tests := []Struct{
		{errs: []error{nil, nil}, want: "nil"},
		{errs: []error{}, want: "nil"},
		{errs: []error{New("new error")}, want: "!nil"},
	}

	for _, tt := range tests {
		result := NewAggregate(tt.errs)

		switch tt.want {
		case "nil":
			if result != nil {
				t.Errorf("NewAggregate([]error) error: want: %v got: %v\n", tt.want, result)
			}
		case "!nil":
			if result == nil {
				t.Errorf("NewAggregate([]error) error: want: %v got: %v\n", tt.want, result)
			}
		}
	}
}

func TestTypeaggregate(t *testing.T) {
	type Struct struct {
		fn         string
		agg        aggregate
		target     error
		wantString string
		wantBool   bool
	}

	same := New("the same error")

	tests := []Struct{
		{
			fn:         "Error",
			agg:        aggregate([]error{New("error1"), New("error2"), NewAggregate([]error{New("error3")})}),
			target:     nil,
			wantString: "[error1; error2; error3]",
			wantBool:   false,
		},
		{
			fn:         "Is",
			agg:        aggregate([]error{New("error1"), same}),
			target:     same,
			wantString: "",
			wantBool:   true,
		},
	}

	for _, tt := range tests {
		switch tt.fn {
		case "Error":
			result := tt.agg.Error()
			if result != tt.wantString {
				t.Errorf("Error(error) error: want: %v got: %v\n", tt.wantString, result)
			}
		case "Is":
			result := tt.agg.Is(tt.target)
			if result != tt.wantBool {
				t.Errorf("Is(error) error: want: %v got: %v\n", tt.wantBool, result)
			}
		}
	}
}
