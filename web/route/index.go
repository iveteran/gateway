package route

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"matrix.works/gateway/conf"
	"matrix.works/gateway/web/middleware"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

func proxyHandler(target string) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(target)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%s [%s] redirect to %s\n", r.URL.Path, r.Method, target)
		requestId := r.Header.Get("X-REQUEST-ID")
		if requestId != "" {
			log.Printf("%s [%s] x-request-id: %s\n", r.URL.Path, r.Method, requestId)
		}

		r.Header.Set("APP-ID", conf.Cfg.Server.AppId)
		r.Header.Set("APP-TOKEN", conf.Cfg.Server.AppToken)

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	}
}

func Setup(routeMap map[string]string) {
	http.Handle("/", middleware.TimerHandler(middleware.Authenticator(http.HandlerFunc(index))))
	for pathPrefix, target := range routeMap {
		http.Handle(pathPrefix, middleware.TimerHandler(middleware.Authenticator(http.HandlerFunc(proxyHandler(target)))))
		pathLen := len(pathPrefix)
		if pathPrefix[pathLen-1] == '/' {
			http.Handle(pathPrefix[:pathLen-1], middleware.TimerHandler(middleware.Authenticator(http.HandlerFunc(proxyHandler(target)))))
		} else {
			http.Handle(pathPrefix+"/", middleware.TimerHandler(middleware.Authenticator(http.HandlerFunc(proxyHandler(target)))))
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to FMX Gateway!\r\n")
}
