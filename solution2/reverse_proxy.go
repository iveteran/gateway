package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jessevdk/go-flags"

	cutils "matrix.works/fmx-common/utils"
	"matrix.works/fmx-gateway/conf"
	"matrix.works/fmx-gateway/datasource"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

func timerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		t1 := time.Now()
		log.Printf("%s cost %v", r.URL.Path, t1.Sub(t0))
	})
}

func authenticatorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlWhiteList := conf.Cfg.Server.UrlWhiteList
		urlPrefixWhiteList := conf.Cfg.Server.UrlPrefixWhiteList
		path := r.URL.Path
		token := r.Header.Get("X-TOKEN")
		uid, _ := strconv.Atoi(r.Header.Get("X-UID"))

		if !cutils.ArrayContains(urlWhiteList, path) &&
			!cutils.ArrayPrefixMatch(urlPrefixWhiteList, path) {
			pass := authenticate(uint32(uid), token)
			if pass {
				log.Printf("%s authorize success", path)
				next.ServeHTTP(w, r)
			} else {
				log.Printf("%s authorize failed", path)
				w.WriteHeader(401)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func authenticate(uid uint32, token string) bool {
	savedToken := getUserToken(uid)
	if savedToken != "" && token != "" && savedToken == token {
		return true
	} else {
		return false
	}
}

func getUserToken(uid uint32) string {
	token := ""
	cache := datasource.GetRedisInstance()
	var key = fmt.Sprintf("user_token:%d", uid)
	results, err := redis.Strings(cache.Do("HMGET", key, "token"))
	if err != nil {
		log.Println("GetUserToken error: ", err)
	}
	if err == nil && len(results) > 0 {
		token = results[0]
	}
	return token
}

func proxyHandler(target string) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(target)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%s redirect %s\n", r.URL.Path, target)

		r.Header.Set("APP-ID", conf.Cfg.Server.AppId)
		r.Header.Set("APP-TOKEN", conf.Cfg.Server.AppToken)

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	}
}

func setupRequestHandlers(routeMap map[string]string) {
	http.Handle("/", timerHandler(authenticatorHandler(http.HandlerFunc(index))))
	for pathPrefix, target := range routeMap {
		http.Handle(pathPrefix, timerHandler(authenticatorHandler(http.HandlerFunc(proxyHandler(target)))))
		pathLen := len(pathPrefix)
		if pathPrefix[pathLen-1] == '/' {
			http.Handle(pathPrefix[:pathLen-1], timerHandler(authenticatorHandler(http.HandlerFunc(proxyHandler(target)))))
		} else {
			http.Handle(pathPrefix+"/", timerHandler(authenticatorHandler(http.HandlerFunc(proxyHandler(target)))))
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!\r\n")
}

type CommandOptions struct {
	ConfigPath  string `short:"c" long:"config" description:"Configure file path"`
	VersionFlag bool   `short:"v" long:"version" description:"Show version"`
}

const (
	AppName    = "fgw"
	AppVersion = "1.0"
	AppBuildNo = "unknown"
	AppOwner   = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
)

func ParseCommandLine(cmdOpts *CommandOptions) {
	parser := flags.NewParser(cmdOpts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Printf("Parse command line: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if cmdOpts.VersionFlag {
		fmt.Printf("Fimatrix(%s) Version: %s, BuildNo: %s, Copyright: %s\n",
			AppName, AppVersion, AppBuildNo, AppOwner)
		os.Exit(0)
	}

	log.Printf("Config file: %s\n", cmdOpts.ConfigPath)

	if _, err := os.Stat(cmdOpts.ConfigPath); os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func main() {
	cmdOpt := &CommandOptions{}
	ParseCommandLine(cmdOpt)
	conf.CreateGlobalConfig(cmdOpt.ConfigPath, nil)

	routeMap := conf.Cfg.RouteTable
	fmt.Printf("route table: %+v\n", routeMap)
	setupRequestHandlers(routeMap)

	addr := fmt.Sprintf("%s:%d", conf.Cfg.Server.ListenAddress, conf.Cfg.Server.ListenPort)
	server := &http.Server{
		Addr: addr,
	}

	log.Printf("Listening on %s...", addr)
	server.ListenAndServe()
}
