package errors

import (
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		coder    DefaultCoder
		err      error
		wantCode int
	}{
		{DefaultCoder{code: 100001, httpStatus: "200", message: "connected", reference: "http://reference.com"}, WithCode(100001, "code error"), 100001},
		{DefaultCoder{code: 100002, httpStatus: "400", message: "invalid token", reference: "http://reference.com"}, WithCode(100002, "code errror"), 100002},
		{DefaultCoder{code: 100003, httpStatus: "500", message: "internal server error", reference: "http://reference.com"}, WithCode(100003, "code error"), 100003},
	}

	for _, tt := range tests {
		Register(&tt.coder)

		if v := ParseCoder(tt.err); v.Code() != tt.wantCode {
			t.Errorf("Register(%v) want:%d got:%d", tt.coder, tt.wantCode, v.Code())
		}
	}
}

func TestMustRegister(t *testing.T) {
	tests := []struct {
		Name        string
		Coder       DefaultCoder
		Err         error
		WantCode    int
		ShouldPanic bool
	}{
		{"regularMustRegister", DefaultCoder{code: 100001, httpStatus: "200", message: "Connected", reference: "http://reference.com"}, WithCode(100001, "code error"), 100001, false},
		{"duplicatedMustRegister", DefaultCoder{code: 100002, httpStatus: "400", message: "Client Error", reference: "http://reference.com"}, WithCode(100002, "code error"), 100002, true},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var Panicked bool
			var PanicValue interface{}

			defer func() {
				if r := recover(); r != nil {
					PanicValue = r
					Panicked = true
				}
			}()

			if tt.ShouldPanic {
				MustRegister(&tt.Coder)
			}

			MustRegister(&tt.Coder)

			if tt.ShouldPanic {
				if !Panicked {
					t.Errorf("MustRegister(%v) expected to panic, but it did't\n", tt.Coder)
				} else {
					if PanicValue != "code already exists" {
						t.Errorf("MustRegister(%v) recover from panic: want:`code already exists` got:%v", tt.Coder, PanicValue)
					}
				}
			} else {
				if Panicked {
					t.Errorf("MustRegister(%v) should not panic, but it did", tt.Coder)
				} else {
					if c := ParseCoder(tt.Err); c.Code() != tt.WantCode {
						t.Errorf("MustRegister(%v) succeed to register but: want:%d got:%d", tt.Coder, tt.WantCode, c.Code())
					}
				}
			}

		})
	}
}
