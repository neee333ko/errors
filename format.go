package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type info struct {
	code  int
	error string
	*stack
}

func (c *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		var jsonMode, withDetail, withCaller bool

		if s.Flag('#') {
			jsonMode = true
		}

		if s.Flag('+') {
			withDetail = true
		}

		if s.Flag('-') {
			withCaller = true
		}

		buffer := bytes.NewBuffer(nil)
		jsonData := make([]map[string]interface{}, 0)
		errlist := list(c)
		sep := ""

		for i, e := range errlist {
			format(i, e, buffer, &jsonData, jsonMode, withDetail, withCaller, sep)

			sep = "; "

			if !withDetail {
				break
			}
		}

		if jsonMode {
			bytes, _ := json.Marshal(jsonData)
			buffer.Write(bytes)
		}

		fmt.Fprint(s, buffer.String())
	default:
		fmt.Fprint(s, c.code, c.cause.Error())
	}
}

func list(e error) []error {
	var ret []error

	if e != nil {
		if w, ok := e.(interface{ Unwrap() error }); ok {
			ret = append(ret, e)
			ret = append(ret, list(w.Unwrap())...)
		} else {
			ret = append(ret, e)
		}
	}

	return ret
}

func format(k int, err error, buffer *bytes.Buffer, jsonData *[]map[string]interface{}, jsonMode,
	withDetail, withCaller bool, sep string) {
	info := buildInfo(err)

	if jsonMode {
		data := map[string]interface{}{}
		if withCaller || withDetail {
			data["code"] = info.code
			data["error"] = info.error

			if info.stack != nil {
				frame := Frame((*info.stack)[0])
				caller := fmt.Sprintf("#%d %s:%d (%s)", k, frame.file(), frame.line(), funcname(frame.name()))
				data["caller"] = caller
			}
		} else {
			data["error"] = info.error
		}

		*jsonData = append(*jsonData, data)
	} else {
		if withCaller || withDetail {
			fmt.Fprintf(buffer, "%s%s", sep, info.error)

			if info.stack != nil {
				frame := Frame((*info.stack)[0])
				fmt.Fprintf(buffer, " - #%d [%s:%d %s] (#%d)",
					k,
					frame.file(),
					frame.line(),
					funcname(frame.name()),
					info.code,
				)
			}
		} else {
			fmt.Fprintf(buffer, "%s%s", sep, info.error)
		}
	}
}

func buildInfo(e error) info {
	var infoItem *info

	switch err := e.(type) {
	case *fundamental:
		infoItem = &info{
			code:  defaultCoder.Code(),
			error: err.Error(),
			stack: err.stack,
		}
	case *withStack:
		infoItem = &info{
			code:  defaultCoder.Code(),
			error: err.Error(),
			stack: err.stack,
		}
	case *withCode:
		infoItem = &info{
			code:  err.code,
			error: err.Error(),
			stack: err.stack,
		}
	default:
		infoItem = &info{
			code:  defaultCoder.Code(),
			error: err.Error(),
		}
	}

	return *infoItem
}
