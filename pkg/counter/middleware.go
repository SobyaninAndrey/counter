package counter

import (
	"net/http"
)

func (c *Counter) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.AddEvent()
		next.ServeHTTP(w, r)
		c.Cancel()
	}
	return http.HandlerFunc(fn)
}
