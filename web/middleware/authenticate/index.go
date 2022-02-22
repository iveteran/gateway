package authenticate

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"

	cutils "matrix.works/fmx-common/utils"
	"matrix.works/fmx-gateway/bootstrap"
	"matrix.works/fmx-gateway/conf"
	"matrix.works/fmx-gateway/datasource"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.Use(func(ctx iris.Context) {
		whiteList := conf.Cfg.Server.UrlWhiteList
		path := ctx.Path()
		uid, _ := strconv.Atoi(ctx.GetHeader("X-UID"))
		token := ctx.GetHeader("X-TOKEN")
		mwxUA := r.Header.Get("MXW-UA")

		if !cutils.ArrayContains(whiteList, path) {
			pass := authenticate(uint32(uid), token, mwxUA)
			if pass {
				ctx.Next()
			} else {
				ctx.StatusCode(401)
				ctx.EndRequest()
			}
		} else {
			ctx.Next()
		}
	})
}

// FIXME(refactor): Call service fmx-user-center or fmx-user-auth to authenticate user
// It's ugly to get user token and match it here
func authenticate(uid uint32, token, client string) bool {
	savedToken := getUserToken(uid, client)
	if savedToken != "" && token != "" && savedToken == token {
		return true
	} else {
		return false
	}
}

func getUserToken(uid uint32, client string) string {
	var key string
	if client != "" {
		key = fmt.Sprintf("user_session@%s:%d", client, uid)
	} else {
		key = fmt.Sprintf("user_session:%d", uid)
	}

	var token string
	cache := datasource.GetRedisDefaultInstance()
	results, err := redis.Strings(cache.Do("HMGET", key, "token"))
	if err != nil {
		log.Println("GetUserToken error: ", err)
	}
	if err == nil && len(results) > 0 {
		token = results[0]
	}
	return token
}
