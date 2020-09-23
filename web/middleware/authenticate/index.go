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
		token := ctx.GetHeader("X-TOKEN")
		uid, _ := strconv.Atoi(ctx.GetHeader("X-UID"))

		if !cutils.ArrayContains(whiteList, path) {
			pass := authenticate(uint32(uid), token)
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
