package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

var (
	UseNumber = true
	DisallowUnknownFields = true
)

type Context struct {
	Req  *http.Request

	// 如果入户直接使用这个
	// 那么他们就绕开了 RespData 和 RespStatusCode 两个
	// 那么部分 middleware 无法运作
	Resp http.ResponseWriter

	// 主要为了 middleware 读写（篡改响应）用的
	RespData []byte
	RespStatusCode int

	PathParams map[string]string

	queryValues url.Values
	// 缓存的数据
	// 这个缓存不存在所谓的失效不一致问题的，因为对于 Web 框架来说，请求收到之后，就是确切无疑，不会改变的
	cacheQueryValues url.Values

	// 匹配的路由
	MatchedRoute string

	tplEngine TemplateEngine
}

func (c *Context) Render(tplName string, data any) error  {
	var err  error
	data, err = c.tplEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
		return err
	}
	c.RespStatusCode = http.StatusOK
	return  nil
}

// BindJSON binds the request body to the given value
func (c *Context) BindJSON(val any, useNumber bool, disallowUnknownFields bool) error {
	if val == nil {
		return errors.New("val is nil")
	}

	if c.Req.Body == nil {
		return errors.New("empty request body")
	}

	decoder := json.NewDecoder(c.Req.Body)

	if useNumber {
		decoder.UseNumber()
	}
	if disallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	// UseNumber => 数字就是用 Number 来表示
	// 否则数字就是用 float64 来表示
	//decoder.UseNumber()
	// DisallowUnknownFields => 禁止未知字段
	// 比如 json 里面有未知字段，就会报错
	//decoder.DisallowUnknownFields()

	return decoder.Decode(val)
}

// FormValue gets the value of the first named param in the query string or form data
func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	/*
	vals, ok := c.Req.Form[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return vals[0], nil

	 */
	return c.Req.FormValue(key), nil
}

// QueryValue gets the value of the first named param in the query string
// Query 和表单比起来，它没有缓存
func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok || len(vals) == 0 {
		return "", errors.New("key not found")
	}
	return vals[0], nil

	// 用户区别不出来是否真的有值，但是值恰好是空字符串
	// 还是没有值
	//return c.Req.URL.Query().Get(key), nil
	//return c.queryValues.Get(key), nil
}

// PathValue gets the value of the first named param in the path
func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}

// PathValue1 gets the value of the first named param in the path
// PathValue 另一种写法，包含了 Int64 的转换
// 使用方法：result.AsInt64()
func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("key not found"),
		}
	}
	return StringValue{
		Value: val,
	}
}

type StringValue struct {
	Value string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.Value, 10, 64)
}

// RespJSON writes the given value as JSON to the response
func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	//c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	c.RespData = data
	c.RespStatusCode = status
	/*
	n, err := c.Resp.Write(data)
	if n != len(data) {
		return errors.New("write data error")
	}

	 */
	return err
}

// SetCookie sets a cookie in the response
func (c *Context) SetCookie(ck *http.Cookie) {
	http.SetCookie(c.Resp, ck)
}

