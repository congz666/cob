package binding

import "net/http"

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type Binding interface {
	Bind(*http.Request, interface{}) error
}

// These implement the Binding interface and can be used to bind the data
var (
	JSON     = jsonBinding{}
	FormPost = formPostBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(contentType string) Binding {
	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEMultipartPOSTForm:
		return FormPost
	default: // case MIMEPOSTForm:
		return FormPost
	}
}
