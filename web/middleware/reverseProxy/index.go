package reverseProxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/kataras/iris"

	"matrix.works/fmx-gateway/bootstrap"
	//"matrix.works/fmx-gateway/conf"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.Use(func(ctx iris.Context) {
		//proxyTable, dftTarget := conf.Cfg.Server.ReverseProxyTable
		path := ctx.Path()
		proxyTable := map[string]string{
			"/query":  "http://127.0.0.1:4020",
			"/market": "http://127.0.0.1:4050",
		}
		dftTarget := "http://127.0.0.1:4020"

		log.Printf("#### path: %s\n", path)
		target := matchProxyTarget(path, proxyTable)
		if target != "" {
			newProxy(target, ctx.ResponseWriter(), ctx.Request())
			ctx.Next()
		} else if dftTarget != "" {
			newProxy(dftTarget, ctx.ResponseWriter(), ctx.Request())
			ctx.Next()
		} else {
			ctx.StatusCode(404)
			ctx.EndRequest()
		}
	})
}

func matchProxyTarget(
	path string,
	proxyTable map[string]string,
) (target string) {
	for matchPath, target := range proxyTable {
		if path == "/" && matchPath == "/" {
			return target
		} else if strings.HasPrefix(path, matchPath) {
			return target
		}
	}
	return target
}

func newProxy(target string, w http.ResponseWriter, r *http.Request) {
	log.Printf("ReverseProxy: redirect: %s -> %s\n", r.URL, target)
	url, err := url.Parse(target)
	if err != nil {
		log.Println(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}
