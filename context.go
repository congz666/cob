package cob

import (
	"cob/binding"
	"encoding/json"
	"fmt"
	"net/http"
)

// H is used for serialization
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

// newContext is the constructor of cob.Context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// Execute the next HandlerFunc
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Packaging error message
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Get the parameters of dynamic routing
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Get form data
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Get route parameter
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Set status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// Set Response-Header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Serialize results with string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Serialize results with JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Output Data
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.GetHeader("Content-Type"))
}

// Get the value of Request-Header
func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

// Bind data to struct
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.ContentType())
	return c.BindData(obj, b)
}

// Call the bind() of the binding instance
func (c *Context) BindData(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Req, obj)
}
