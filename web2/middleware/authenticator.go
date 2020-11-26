package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"

	cutils "matrix.works/fmx-common/utils"
	"matrix.works/fmx-gateway/conf"
	"matrix.works/fmx-gateway/datasource"
)

func Authenticator(next http.Handler) http.Handler {
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
