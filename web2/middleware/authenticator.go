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
		urlUserAccessList := conf.Cfg.Server.UrlUserAccessList
		path := r.URL.Path
		token := r.Header.Get("X-TOKEN")
		uid, _ := strconv.Atoi(r.Header.Get("X-UID"))

		if CheckUrlDontNeedAuthenticateForUser(uint32(uid), path,
			urlWhiteList, urlPrefixWhiteList, urlUserAccessList) {
			next.ServeHTTP(w, r)
		} else {
			pass := authenticate(uint32(uid), token)
			if pass {
				log.Printf("%s authorize success", path)
				next.ServeHTTP(w, r)
			} else {
				log.Printf("%s authorize failed, uid: %d, token: %s", path, uid, token)
				w.WriteHeader(401)
			}
		}
	})
}

func CheckUrlDontNeedAuthenticateForUser(
	uid uint32, path string,
	whiteList, urlPrefixWhiteList, urlUserAccessList []string,
) bool {
	// 1) path在白名单的可以直接访问，2) 访客可以访问非用户关联接口
	return cutils.ArrayContains(whiteList, path) ||
		cutils.ArrayPrefixMatch(urlPrefixWhiteList, path) ||
		(uid == conf.Cfg.Server.GuestUserId && !cutils.ArrayPrefixMatch(urlUserAccessList, path))
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
