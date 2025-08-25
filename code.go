package errors

import (
	"fmt"
	"sync"
)

var defaultCoder DefaultCoder = DefaultCoder{code: 100000, httpStatus: "200", message: "", reference: ""}

type Coder interface {
	Code() int
	HttpStatus() string
	Message() string
	Reference() string
}

type DefaultCoder struct {
	code       int
	httpStatus string
	message    string
	reference  string
}

func (c *DefaultCoder) Code() int {
	return c.code
}

func (c *DefaultCoder) HttpStatus() string {
	return c.httpStatus
}

func (c *DefaultCoder) Message() string {
	return c.message
}

func (c *DefaultCoder) Reference() string {
	return c.reference
}

var (
	pool = make(map[int]Coder)
	mu   sync.RWMutex
)

func Register(c Coder) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := pool[c.Code()]; !ok {
		pool[c.Code()] = c
		return nil
	}

	return fmt.Errorf("code already exists")
}

func MustRegister(c Coder) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := pool[c.Code()]; !ok {
		pool[c.Code()] = c
		return
	}

	panic(fmt.Errorf("code already exists"))
}

func GetCoder(code int) Coder {
	mu.RLock()
	defer mu.RUnlock()

	if v, ok := pool[code]; ok {
		return v
	}

	return &defaultCoder
}
