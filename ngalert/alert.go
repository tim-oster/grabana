package alert

import (
	"fmt"

	"github.com/K-Phoen/grabana/ngalert/expr"
	"github.com/K-Phoen/grabana/ngalert/query"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure an alert.
type Option func(alert *Alert)

type Alert struct {
	Builder *sdk.NgAlert

	// RefPanelTitle is for internal use for finding the panel the alert was created on.
	RefPanelTitle *string
}

// New creates a new alert.
func New(name string, options ...Option) *Alert {
	alert := &Alert{
		Builder: &sdk.NgAlert{Title: name},
	}
	for _, opt := range append(defaults(), options...) {
		opt(alert)
	}
	if alert.Builder.Condition == "" {
		panic("no query or expression is marked as alert condition")
	}
	return alert
}

func defaults() []Option {
	return []Option{
		For("5m"),
		OnNoData(sdk.NoDataStateNoData),
		OnExecutionError(sdk.ExecErrorStateAlerting),
	}
}

func (alert *Alert) HookDatasource(datasourcesMap map[string]string) error {
	for i, data := range alert.Builder.Data {
		if data.DatasourceUid == "__expr__" {
			continue
		}

		datasourceUid := datasourcesMap[data.DatasourceUid]
		if datasourceUid == "" {
			return fmt.Errorf("could not infer datasource UID from its name: %s", data.DatasourceUid)
		}
		alert.Builder.Data[i].DatasourceUid = datasourceUid
	}

	return nil
}

func (alert *Alert) HookDashboardUID(uid string) {
	Annotate("__dashboardUid__", uid)(alert)
}

func (alert *Alert) HookPanelID(id string) {
	Annotate("__panelId__", id)(alert)
}

// FolderUID defines the uid of the folder the alert belongs to.
func FolderUID(uid string) Option {
	return func(alert *Alert) {
		alert.Builder.FolderUID = uid
	}
}

// RuleGroup defines the rule group the alert belongs to.
func RuleGroup(name string) Option {
	return func(alert *Alert) {
		alert.Builder.RuleGroup = name
	}
}

// OnExecutionError defines the behavior on execution error.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#execution-errors-or-timeouts
func OnExecutionError(state sdk.ExecErrorState) Option {
	return func(alert *Alert) {
		alert.Builder.ExecErrState = state
	}
}

// OnNoData defines the behavior when the query returns no data.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#no-data-null-values
func OnNoData(state sdk.NoDataState) Option {
	return func(alert *Alert) {
		alert.Builder.NoDataState = state
	}
}

// For sets the time interval during which a query violating the threshold
// before the alert being actually triggered.
// See https://grafana.com/docs/grafana/latest/alerting/rules/#for
func For(duration string) Option {
	return func(alert *Alert) {
		alert.Builder.For = duration
	}
}

// Annotate adds a new annotation to the alert rule.
func Annotate(key, value string) Option {
	return func(alert *Alert) {
		if alert.Builder.Annotations == nil {
			alert.Builder.Annotations = map[string]string{}
		}
		alert.Builder.Annotations[key] = value
	}
}

// Summary sets the summary associated to the alert.
func Summary(content string) Option {
	return Annotate("summary", content)
}

// Description sets the description associated to the alert.
func Description(content string) Option {
	return Annotate("description", content)
}

// Runbook sets the runbook URL associated to the alert.
func Runbook(url string) Option {
	return Annotate("runbook_url", url)
}

// Label adds a new label that will be forwarded to the notifications
// channels when the alert will tbe triggered or used to route the alert.
func Label(key, value string) Option {
	return func(alert *Alert) {
		if alert.Builder.Labels == nil {
			alert.Builder.Labels = map[string]string{}
		}
		alert.Builder.Labels[key] = value
	}
}

func Query(refId string, opts ...query.Option) Option {
	return func(alert *Alert) {
		q := query.New(refId, opts...)
		alert.Builder.Data = append(alert.Builder.Data, *q.Builder)

		if q.IsAlertCondition {
			alert.Builder.Condition = refId
		}
	}
}

func Expr(refId string, opts ...expr.Option) Option {
	return func(alert *Alert) {
		e := expr.New(refId, opts...)
		alert.Builder.Data = append(alert.Builder.Data, *e.Builder)

		if e.IsAlertCondition {
			alert.Builder.Condition = refId
		}
	}
}
