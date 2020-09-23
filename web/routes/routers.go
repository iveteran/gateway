package routes

import (
	"github.com/kataras/iris"

	"matrix.works/fmx-gateway/bootstrap"
)

func Configure(b *bootstrap.Bootstrapper) {
	println("#### route configure")
	b.Any("/", index)
	b.Any("/query/", index)
	b.Any("/market/", index)
	// TODO: 实现通用的router, 将所有的path映射到index处理, 实现类似: b.Any("/*", index)
}

func index(ctx iris.Context) {
	ctx.Application().Logger().Infof("#### routing path: %s", ctx.Path())
	ctx.Next()
}
