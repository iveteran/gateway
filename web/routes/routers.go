package routes

import (
	"matrix.works/fmx-common/iris"
	"matrix.works/fmx-gateway/bootstrap"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.Any("/", index)
	b.Any("/query/", index)
	b.Any("/market/", index)
	b.Any("/user/", index)
	// TODO: 实现通用的router, 将所有的path映射到index处理, 实现类似: b.Any("/*", index)
}

func index(ctx iris.Context) {
	ctx.Application().Logger().Infof("#### routing path: %s", ctx.Path())
	ctx.Next()
}
