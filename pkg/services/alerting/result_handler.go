package alerting

import (
	"time"

	"github.com/sound-of-destiny/qlsc_zhxf/pkg/bus"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/components/simplejson"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/log"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/metrics"
	m "github.com/sound-of-destiny/qlsc_zhxf/pkg/models"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/services/annotations"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/services/rendering"
)

type ResultHandler interface {
	Handle(evalContext *EvalContext) error
}

type DefaultResultHandler struct {
	notifier NotificationService
	log      log.Logger
}

func NewResultHandler(renderService rendering.Service) *DefaultResultHandler {
	return &DefaultResultHandler{
		log:      log.New("alerting.resultHandler"),
		notifier: NewNotificationService(renderService),
	}
}

func (handler *DefaultResultHandler) Handle(evalContext *EvalContext) error {
	executionError := ""
	annotationData := simplejson.New()

	if len(evalContext.EvalMatches) > 0 {
		annotationData.Set("evalMatches", simplejson.NewFromAny(evalContext.EvalMatches))
	}

	if evalContext.Error != nil {
		executionError = evalContext.Error.Error()
		annotationData.Set("error", executionError)
	} else if evalContext.NoDataFound {
		annotationData.Set("noData", true)
	}

	metrics.M_Alerting_Result_State.WithLabelValues(string(evalContext.Rule.State)).Inc()
	if evalContext.ShouldUpdateAlertState() {
		handler.log.Info("New state change", "alertId", evalContext.Rule.Id, "newState", evalContext.Rule.State, "prev state", evalContext.PrevAlertState)

		cmd := &m.SetAlertStateCommand{
			AlertId:  evalContext.Rule.Id,
			OrgId:    evalContext.Rule.OrgId,
			State:    evalContext.Rule.State,
			Error:    executionError,
			EvalData: annotationData,
		}

		if err := bus.Dispatch(cmd); err != nil {
			if err == m.ErrCannotChangeStateOnPausedAlert {
				handler.log.Error("Cannot change state on alert that's paused", "error", err)
				return err
			}

			if err == m.ErrRequiresNewState {
				handler.log.Info("Alert already updated")
				return nil
			}

			handler.log.Error("Failed to save state", "error", err)
		}

		// save annotation
		item := annotations.Item{
			OrgId:       evalContext.Rule.OrgId,
			DashboardId: evalContext.Rule.DashboardId,
			PanelId:     evalContext.Rule.PanelId,
			AlertId:     evalContext.Rule.Id,
			Text:        "",
			NewState:    string(evalContext.Rule.State),
			PrevState:   string(evalContext.PrevAlertState),
			Epoch:       time.Now().UnixNano() / int64(time.Millisecond),
			Data:        annotationData,
		}

		annotationRepo := annotations.GetRepository()
		if err := annotationRepo.Save(&item); err != nil {
			handler.log.Error("Failed to save annotation for new alert state", "error", err)
		}
	}

	if evalContext.Rule.State == m.AlertStateOK && evalContext.PrevAlertState != m.AlertStateOK {
		for _, notifierId := range evalContext.Rule.Notifications {
			cmd := &m.CleanNotificationJournalCommand{
				AlertId:    evalContext.Rule.Id,
				NotifierId: notifierId,
				OrgId:      evalContext.Rule.OrgId,
			}
			if err := bus.DispatchCtx(evalContext.Ctx, cmd); err != nil {
				handler.log.Error("Failed to clean up old notification records", "notifier", notifierId, "alert", evalContext.Rule.Id, "Error", err)
			}
		}
	}

	handler.notifier.SendIfNeeded(evalContext)
	return nil
}
