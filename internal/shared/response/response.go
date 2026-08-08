package response

import (
	"encoding/json"
	"net/http"
)

type Response map[string]any

func New() Response {
	return make(Response)
}

func (r Response) Response(msg any) Response {
	r["message"] = msg
	return r
}

func (r Response) InternalServerError(w http.ResponseWriter) {
	r.Error("Internal Server Error", ErrorCodeUnknown).Send(w, 500)
}

func (r Response) Send(w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(r)
}
