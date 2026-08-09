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
	r["response"] = msg
	return r
}

func (r Response) InternalServerError(w http.ResponseWriter) {
	r.Error("Internal Server Error", ErrorCodeUnknown).Send(w, http.StatusBadRequest)
}

func (r Response) BadRequest(w http.ResponseWriter) {
	r.Error("Bad Request", ErrorCodeBadRequest).Send(w, http.StatusBadRequest)
}

func (r Response) Send(w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(r)
}
