package alerting

import (
	"context"
	"fmt"

	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/components/simplejson"
	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
)

type AlertTestCommand struct {
	Dashboard *simplejson.Json
	PanelId   int64
	OrgId     int64

	Result *EvalContext
}

func init() {
	bus.AddHandler("alerting", handleAlertTestCommand)
}

func handleAlertTestCommand(cmd *AlertTestCommand) error {

	dash := m.NewDashboardFromJson(cmd.Dashboard)

	extractor := NewDashAlertExtractor(dash, cmd.OrgId)
	alerts, err := extractor.GetAlerts()
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		if alert.PanelId == cmd.PanelId {
			rule, err := NewRuleFromDBAlert(alert)
			if err != nil {
				return err
			}

			cmd.Result = testAlertRule(rule)
			return nil
		}
	}

	return fmt.Errorf("Could not find alert with panel id %d", cmd.PanelId)
}

func testAlertRule(rule *Rule) *EvalContext {
	handler := NewEvalHandler()

	context := NewEvalContext(context.Background(), rule)
	context.IsTestRun = true

	handler.Eval(context)
	context.Rule.State = context.GetNewState()

	return context
}