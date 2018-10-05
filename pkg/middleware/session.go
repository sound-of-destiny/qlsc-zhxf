package middleware

import (
	ms "github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/services/session"
)

func Sessioner(options *ms.Options, sessionConnMaxLifetime int64) macaron.Handler {
	session.Init(options, sessionConnMaxLifetime)

	return func(ctx *m.ReqContext) {
		ctx.Next()

		if err := ctx.Session.Release(); err != nil {
			panic("session(release): " + err.Error())
		}
	}
}
