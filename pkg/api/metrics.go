package api

import (
	"context"

	"github.com/sound-of-destiny/qlsc_zhxf/pkg/api/dtos"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/components/simplejson"
	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/tsdb"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/tsdb/testdata"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/util"
)

// POST /api/tsdb/query
func (hs *HTTPServer) QueryMetrics(c *m.ReqContext, reqDto dtos.MetricRequest) Response {
	timeRange := tsdb.NewTimeRange(reqDto.From, reqDto.To)

	if len(reqDto.Queries) == 0 {
		return Error(400, "No queries found in query", nil)
	}

	datasourceId, err := reqDto.Queries[0].Get("datasourceId").Int64()
	if err != nil {
		return Error(400, "Query missing datasourceId", nil)
	}

	ds, err := hs.getDatasourceFromCache(datasourceId, c)
	if err != nil {
		return Error(500, "Unable to load datasource meta data", err)
	}

	request := &tsdb.TsdbQuery{TimeRange: timeRange}

	for _, query := range reqDto.Queries {
		request.Queries = append(request.Queries, &tsdb.Query{
			RefId:         query.Get("refId").MustString("A"),
			MaxDataPoints: query.Get("maxDataPoints").MustInt64(100),
			IntervalMs:    query.Get("intervalMs").MustInt64(1000),
			Model:         query,
			DataSource:    ds,
		})
	}

	resp, err := tsdb.HandleRequest(c.Req.Context(), ds, request)
	if err != nil {
		return Error(500, "Metric request error", err)
	}

	statusCode := 200
	for _, res := range resp.Results {
		if res.Error != nil {
			res.ErrorString = res.Error.Error()
			resp.Message = res.ErrorString
			statusCode = 400
		}
	}

	return JSON(statusCode, &resp)
}

// GET /api/tsdb/testdata/scenarios
func GetTestDataScenarios(c *m.ReqContext) Response {
	result := make([]interface{}, 0)

	for _, scenario := range testdata.ScenarioRegistry {
		result = append(result, map[string]interface{}{
			"id":          scenario.Id,
			"name":        scenario.Name,
			"description": scenario.Description,
			"stringInput": scenario.StringInput,
		})
	}

	return JSON(200, &result)
}

// Generates a index out of range error
func GenerateError(c *m.ReqContext) Response {
	var array []string
	return JSON(200, array[20])
}

// GET /api/tsdb/testdata/gensql
func GenerateSQLTestData(c *m.ReqContext) Response {
	if err := bus.Dispatch(&m.InsertSqlTestDataCommand{}); err != nil {
		return Error(500, "Failed to insert test data", err)
	}

	return JSON(200, &util.DynMap{"message": "OK"})
}

// GET /api/tsdb/testdata/random-walk
func GetTestDataRandomWalk(c *m.ReqContext) Response {
	from := c.Query("from")
	to := c.Query("to")
	intervalMs := c.QueryInt64("intervalMs")

	timeRange := tsdb.NewTimeRange(from, to)
	request := &tsdb.TsdbQuery{TimeRange: timeRange}

	dsInfo := &m.DataSource{Type: "testdata"}
	request.Queries = append(request.Queries, &tsdb.Query{
		RefId:      "A",
		IntervalMs: intervalMs,
		Model: simplejson.NewFromAny(&util.DynMap{
			"scenario": "random_walk",
		}),
		DataSource: dsInfo,
	})

	resp, err := tsdb.HandleRequest(context.Background(), dsInfo, request)
	if err != nil {
		return Error(500, "Metric request error", err)
	}

	return JSON(200, &resp)
}
