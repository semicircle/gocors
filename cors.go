package gocors

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Cors struct {
	allowOrigin      string
	allowMethods     map[string]bool // comma-delimited
	allowHeaders     map[string]bool
	allowCredentials bool
	exposeHeaders    string
	maxAge           int
	userHandler      http.Handler
}

func New() *Cors {
	c := new(Cors)
	c.allowOrigin = "*"
	c.allowCredentials = false
	c.allowHeaders = map[string]bool{"origin": true}
	c.allowMethods = map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "UPDATE": true}
	c.maxAge = 1500
	return c
}

func arrayToSet(array []string) map[string]bool {
	ret := make(map[string]bool, len(array))
	for _, v := range array {
		ret[v] = true
	}
	return ret
}

func setToArray(set map[string]bool) []string {
	ret := make([]string, 0, len(set))
	for k, _ := range set {
		ret = append(ret, k)
	}
	return ret
}

func (cors *Cors) SetAllowOrigin(origin string) *Cors {
	cors.allowOrigin = origin
	return cors
}

func (cors *Cors) AllowOrigin() string {
	return cors.allowOrigin
}

func (cors *Cors) SetMaxAge(maxAge int) *Cors {
	cors.maxAge = maxAge
	return cors
}

func (cors *Cors) MaxAge() int {
	return cors.maxAge
}

func (cors *Cors) SetAllowMethods(methods []string) *Cors {
	cors.allowMethods = arrayToSet(methods)
	return cors
}

func (cors *Cors) AllowMethods() []string {
	return setToArray(cors.allowMethods)
}

func (cors *Cors) SetAllowHeaders(headers []string) *Cors {
	cors.allowHeaders = arrayToSet(headers)
	return cors
}

func (cors *Cors) AllowHeaders() []string {
	return setToArray(cors.allowHeaders)
}

func (cors *Cors) SetExposeHeaders(headers string) *Cors {
	cors.exposeHeaders = headers
	return cors
}

func (cors *Cors) ExposeHeaders() string {
	return cors.exposeHeaders
}

func (cors *Cors) SetAllowCredentials(allow bool) *Cors {
	cors.allowCredentials = allow
	return cors
}

func (cors *Cors) AllowCredentials() bool {
	return cors.allowCredentials
}

func (cors *Cors) Handler(h http.Handler) http.Handler {
	//cors = &Cors{cors.allowOrigin, cors.allowMethods, cors.allowHeaders, cors.allowCredentials, cors.exposeHeaders, cors.maxAge, h}
	cors.userHandler = h
	return cors
}

func (cors *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin == "" {
		//cors.corsNotValid(w, r)
		cors.userHandler.ServeHTTP(w, r)
		return
	} else if r.Method != "OPTIONS" {
		//actual request.
		cors.actualRequest(w, r)
		return
	} else if acrm := r.Header.Get("Access-Control-Request-Method"); acrm == "" {
		//actual request.
		cors.actualRequest(w, r)
		return
	} else {
		//preflight request.
		cors.preflightRequest(w, r)
		return
	}
}

func (cors *Cors) preflightRequest(w http.ResponseWriter, r *http.Request) {
	acrm := r.Header.Get("Access-Control-Request-Method")
	//log.Printf("ACRM: %s", acrm)
	if acrm == "" {
		cors.corsNotValid(w, r)
		log.Printf("Access-Control-Request-Method: is empty")
		return
	}
	methods := strings.Split(strings.TrimSpace(acrm), ",")
	for _, m := range methods {
		m = strings.TrimSpace(m)
		if _, ok := cors.allowMethods[m]; !ok {
			cors.corsNotValid(w, r)
			log.Printf("Access-Control-Request-Method: %s is not supported", m)
			return
		}
	}
	acrh := r.Header.Get("Access-Control-Request-Headers")
	//log.Printf("ACRH: %s", acrh)
	if acrh != "" {
		headers := strings.Split(strings.TrimSpace(acrh), ",")
		for _, h := range headers {
			h = strings.TrimSpace(h)
			if _, ok := cors.allowHeaders[h]; !ok {
				cors.corsNotValid(w, r)
				log.Printf("Access-Control-Request-Headers: `%s` is not supported \n", h)
				log.Printf("cors.allowHeaders: `%v`\n", cors.allowHeaders)
				return
			}
		}
	}

	//TODO: make this faster.
	w.Header().Add("Access-Control-Allow-Methods", strings.Join(cors.AllowMethods(), ","))
	w.Header().Add("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders(), ","))
	w.Header().Add("Access-Control-Max-Age", strconv.Itoa(cors.maxAge))
	cors.setAllowOrigin(w, r)
	cors.setAllowCookies(w, r)
	return
}

func (cors *Cors) corsNotValid(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "CORS Request Invalid")
	return
}

func (cors *Cors) actualRequest(w http.ResponseWriter, r *http.Request) {
	if cors.exposeHeaders != "" {
		w.Header().Add("Access-Control-Expose-Headers", cors.exposeHeaders)
	}
	cors.setAllowOrigin(w, r)
	cors.setAllowCookies(w, r)
	cors.userHandler.ServeHTTP(w, r)
	return
}

func (cors *Cors) setAllowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", cors.allowOrigin)
	return
}

func (cors *Cors) setAllowCookies(w http.ResponseWriter, r *http.Request) {
	if cors.allowCredentials {
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}
