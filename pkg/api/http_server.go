package api

import (
	"context"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/api/live"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/api/routing"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/log"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/services/rendering"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/setting"
	"gopkg.in/macaron.v1"
	"net/http"
)

type HTTPServer struct {
	log           log.Logger
	macaron       *macaron.Macaron
	context       context.Context
	streamManager *live.StreamManager
	//cache         *gocache.Cache
	httpSrv *http.Server

	RouteRegister routing.RouteRegister `inject:""`
	Bus           bus.Bus               `inject:""`
	RenderService rendering.Service     `inject:""`
	Cfg           *setting.Cfg          `inject:""`
}
