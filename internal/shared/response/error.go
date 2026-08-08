package response

type ErrorField struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r Response) Error(msg, code string) Response {
	errs, ok := r["errors"]
	if !ok {
		errs = make([]ErrorField, 0)
	}
	errs = append(errs.([]ErrorField), ErrorField{
		Message: msg,
		Code:    code,
	})
	r["errors"] = errs
	return r
}
