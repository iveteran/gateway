package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"

	"matrix.works/fmx-common/datasource"
	cutils "matrix.works/fmx-common/utils"
	"matrix.works/fmx-gateway/conf"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlWhiteList := conf.Cfg.UrlPermission.UrlWhiteList
		urlPrefixWhiteList := conf.Cfg.UrlPermission.UrlPrefixWhiteList
		urlUserAccessList := conf.Cfg.UrlPermission.UrlUserAccessList
		path := r.URL.Path
		uid, _ := strconv.Atoi(r.Header.Get("X-UID"))
		token := r.Header.Get("X-TOKEN")
		mwxUA := r.Header.Get("MWX-UA")
		machineId := r.Header.Get("MACHINE-ID")

		if CheckUrlDontNeedAuthenticateForUser(uint32(uid), path,
			urlWhiteList, urlPrefixWhiteList, urlUserAccessList) {
			next.ServeHTTP(w, r)
		} else {
			pass := authenticate(uint32(uid), token, mwxUA, machineId)
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
		(uid == conf.Cfg.Misc.GuestUserId && !cutils.ArrayPrefixMatch(urlUserAccessList, path))
}

// FIXME(refactor): Call service fmx-user-center or fmx-user-auth to authenticate user
// It's ugly to get user token and match it here
func authenticate(uid uint32, token, client, machineId string) bool {
	savedToken, savedMachineId := getUserSession(uid, client)
	if (savedToken != "" && token != "" && savedToken == token) &&
		(savedMachineId == "" || savedMachineId == machineId) {
		return true
	} else {
		return false
	}
}

func getUserSession(uid uint32, client string) (string, string) {
	var key string
	if client != "" {
		key = fmt.Sprintf("user_session@%s:%d", client, uid)
	} else {
		key = fmt.Sprintf("user_session:%d", uid)
	}

	var token, machineId string
	cache := datasource.GetRedisDefaultInstance()
	results, err := redis.Strings(cache.Do("HMGET", key, "token", "machine_id"))
	if err != nil {
		log.Println("GetUserSession error: ", err)
	}
	if err == nil && len(results) > 0 {
		token = results[0]
		machineId = results[1]
	}
	return token, machineId
}
