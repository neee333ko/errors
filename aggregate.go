package errors

import (
	"bytes"
	"fmt"
)

type MessageCountMap map[string]int

type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

func NewAggregate(errs []error) Aggregate {
	if errs == nil {
		return nil
	}

	e := make([]error, 0)

	for _, err := range errs {
		if err != nil {
			e = append(e, err)
		}
	}

	if len(e) == 0 {
		return nil
	}

	return aggregate(e)

}

type aggregate []error

func (agg aggregate) Error() string {
	s := NewString()

	agg.visit(func(e error) bool {
		message := e.Error()
		s.Insert(message)
		return false
	})

	b := bytes.NewBuffer(nil)

	ss := make([]string, 0)
	for key := range s {
		ss = append(ss, key)
	}

	for i, v := range ss {
		if i == len(ss)-1 {
			b.WriteString(v)
		} else {
			b.WriteString(v + "; ")
		}
	}

	return "[" + b.String() + "]"
}

func (agg aggregate) Errors() []error {
	return []error(agg)
}

func (agg aggregate) Is(target error) bool {
	return agg.visit(func(e error) bool {
		return Is(e, target)
	})
}

func (agg aggregate) visit(f func(error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate:
			if err.visit(f) {
				return true
			}
		case Aggregate:
			for _, e := range err.Errors() {
				if f(e) {
					return true
				}
			}
		default:
			if f(err) {
				return true
			}
		}
	}

	return false
}

type Matcher func(error) bool

func FilterOut(err error, matchers ...Matcher) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(Aggregate); ok {
		return NewAggregate(filter(err, matchers...))
	}

	for _, matcher := range matchers {
		if matcher(err) {
			return nil
		}
	}

	return err
}

func filter(agg Aggregate, matchers ...Matcher) []error {
	errs := make([]error, 0, len(agg.Errors()))

	for _, err := range agg.Errors() {
		if err = FilterOut(err, matchers...); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func Flatten(agg Aggregate) Aggregate {
	new := make([]error, 0, len(agg.Errors()))

	for _, err := range agg.Errors() {
		if e, ok := err.(Aggregate); ok {
			new = append(new, Flatten(e).Errors()...)
		} else {
			new = append(new, err)
		}
	}

	return NewAggregate(new)
}

func MessageCountMapToAggregate(theMap MessageCountMap) Aggregate {
	errs := make([]error, 0, len(theMap))

	for key, value := range theMap {
		if value >= 1 {
			errs = append(errs, fmt.Errorf(key))
		}
	}

	return NewAggregate(errs)
}

func Reduce(err error) error {
	if agg, ok := err.(Aggregate); ok {
		errs := agg.Errors()
		switch len(errs) {
		case 0:
			return nil
		case 1:
			return errs[0]
		}
	}

	return err
}

func AggregateGoRoutines(fns ...func() error) Aggregate {
	errs := make([]error, 0, len(fns))

	errsChan := make(chan error, len(fns))

	for _, fn := range fns {
		go func(f func() error) {
			errsChan <- f()
		}(fn)
	}

	for err := range errsChan {
		errs = append(errs, err)
	}

	return NewAggregate(errs)
}
