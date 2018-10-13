package sqlstore

import (
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}
