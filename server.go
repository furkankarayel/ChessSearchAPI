package engine

import (
	"encoding/json"
	"engine/middleware"
	"net/http"
	"path"
	"strings"
)

type Route struct {
	WithLogger bool
	Handler    http.Handler
}

type Server struct {
	Routes map[string]*Route
}

func New(routes map[string]*Route) *Server {
	return &Server{Routes: routes}
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	route, ok := s.Routes[head]
	if !ok {
		Respond(w, r, http.StatusNotFound, "root route not found")
		return
	}

	next := route.Handler

	if route.WithLogger {
		next = middleware.Logger(next)
	}
	next.ServeHTTP(w, r)
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {
	if e, ok := data.(error); ok {
		var tmp = new(struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		})
		tmp.Status = "error"
		tmp.Error = e.Error()
		data = tmp
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	middleware.LogRequest(r, status)

	return nil
}
