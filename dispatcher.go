package ex

import (
	"reflect"
	"strings"
)

type Dispatcher struct {
	ctrl map[string]any
}

func newDispatcher() *Dispatcher {
	return &Dispatcher{
		ctrl: make(map[string]any),
	}
}

func (d *Dispatcher) RegisterDispatcher(name string, ctrl any) {
	d.ctrl[name] = ctrl
}

func (d *Dispatcher) Dispatch(ctx *Context) bool {
	path := strings.Trim(ctx.Path, "/")
	if path == "" {
		return false
	}

	parts := strings.Split(path, "/")
	module := parts[0]
	ctl, ok := d.ctrl[module]
	if !ok {
		return false
	}

	v := reflect.ValueOf(ctl)

	methodName := strings.Title(strings.ToLower(ctx.Method))
	if len(parts) > 1 {
		methodName = strings.Title(parts[1])
	}
	method := v.MethodByName(methodName)
	if !method.IsValid() {
		return false
	}

	method.Call([]reflect.Value{reflect.ValueOf(ctx)})

	return true
}
