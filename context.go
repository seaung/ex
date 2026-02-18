package ex

/*
 *  用于封装请求上下文
 */
import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

type HandlerFunc func(*Context)

// 请求上下文结构体
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	handlers   []HandlerFunc
	index      int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (ctx *Context) Abort() {
	ctx.index = len(ctx.handlers)
}

func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < len(ctx.handlers) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

// 响应http string
func (ctx *Context) String(code int, msg string) {
	ctx.Writer.WriteHeader(code)
	ctx.Writer.Write([]byte(msg))
}

// 响应http json
func (ctx *Context) Json(code int, obj any) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Status(code)
	enCoder := json.NewEncoder(ctx.Writer)
	if err := enCoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

// 获取URL参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) ShouldBindJson(obj any) error {
	if ctx.Req.Body == nil {
		return http.ErrBodyNotAllowed
	}
	decoder := json.NewDecoder(ctx.Req.Body)
	return decoder.Decode(obj)
}

func (ctx *Context) ShouldBindQuery(obj any) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return errors.New("ojbect must be a non-nil pointer")
	}
	val = val.Elem()
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		queryTag := field.Tag.Get("query")
		if queryTag == "" {
			queryTag = field.Name
		}
		values, ok := ctx.Req.URL.Query()[queryTag]
		if !ok || len(values) == 0 {
			continue
		}

		fv := val.Field(i)
		if !val.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(values[0])
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return err
			}
			fv.SetInt(intVal)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintVal, err := strconv.ParseUint(values[0], 10, 64)
			if err != nil {
				return err
			}
			fv.SetUint(uintVal)
		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				return err
			}
			fv.SetFloat(floatVal)
		case reflect.Bool:
			bVal, err := strconv.ParseBool(values[0])
			if err != nil {
				return err
			}
			fv.SetBool(bVal)
		default:
			return errors.New("nosupport type")
		}
	}
	return nil
}
