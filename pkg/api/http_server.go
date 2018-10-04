package api

import (
	"context"
	"github.com/grafana/grafana/pkg/api/live"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/services/rendering"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/log"
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
	httpSrv       *http.Server

	RouteRegister routing.RouteRegister `inject:""`
	Bus           bus.Bus               `inject:""`
	RenderService rendering.Service     `inject:""`
	Cfg           *setting.Cfg          `inject:""`
}