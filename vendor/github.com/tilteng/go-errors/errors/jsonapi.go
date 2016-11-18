package errors

type JSONAPIErrorResponse struct {
	Errors JSONAPIErrors `json:"errors"`
}

type JSONAPIErrors []*JSONAPIError

type JSONAPIErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

type JSONAPIErrorMeta map[string]interface{}

type JSONAPIErrorLinks struct {
	About string `json:"about"`
}

type JSONAPIError struct {
	ID     string              `json:"id"`
	Status int                 `json:"status,string"`
	Links  *JSONAPIErrorLinks  `json:"links,omitempty"`
	Code   string              `json:"code,omitempty"`
	Title  string              `json:"title,omitempty"`
	Detail string              `json:"detail,omitempty"`
	Source *JSONAPIErrorSource `json:"source,omitempty"`
	Meta   JSONAPIErrorMeta    `json:"meta,omitempty"`
}

func (self *JSONAPIErrors) AddError(err *JSONAPIError) {
	*self = append(*self, err)
}
