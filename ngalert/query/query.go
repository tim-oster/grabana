package query

import (
	"github.com/K-Phoen/sdk"
	"time"
)

type Option func(query *Query)

type Query struct {
	Builder          *sdk.NgAlertQuery
	IsAlertCondition bool
}

func New(refId string, options ...Option) *Query {
	const invalid = "INVALID"

	query := &Query{
		Builder: &sdk.NgAlertQuery{
			RefId: refId,
			RelativeTimeRange: sdk.RelativeTimeRange{
				FromSeconds: 600,
				ToSeconds:   0,
			},
			DatasourceUid: "",
			Model: sdk.NgAlertQueryModel{
				NgAlertQueryModelQuery: &sdk.NgAlertQueryModelQuery{
					EditorMode:    "code",
					Expr:          invalid,
					Instant:       false,
					IntervalMs:    sdk.DefaultIntervalMs,
					LegendFormat:  "__auto",
					MaxDataPoints: sdk.DefaultMaxDataPoints,
					Range:         true,
					RefId:         refId,
				},
			},
		},
	}
	for _, opt := range options {
		opt(query)
	}

	if query.Builder.Model.NgAlertQueryModelQuery.Expr == invalid {
		panic("invalid query expression")
	}

	return query
}

func AlertCondition() Option {
	return func(query *Query) {
		query.IsAlertCondition = true
	}
}

func TimeRange(from, to time.Duration) Option {
	return func(query *Query) {
		query.Builder.RelativeTimeRange.FromSeconds = int(from.Seconds())
		query.Builder.RelativeTimeRange.ToSeconds = int(to.Seconds())
	}
}

func Datasource(uid string) Option {
	return func(query *Query) {
		query.Builder.DatasourceUid = uid
	}
}

func Expr(expr string) Option {
	return func(query *Query) {
		query.Builder.Model.NgAlertQueryModelQuery.Expr = expr
	}
}

func Instant() Option {
	return func(query *Query) {
		query.Builder.Model.NgAlertQueryModelQuery.Instant = true
		query.Builder.Model.NgAlertQueryModelQuery.Range = false
	}
}

func Range() Option {
	return func(query *Query) {
		query.Builder.Model.NgAlertQueryModelQuery.Instant = false
		query.Builder.Model.NgAlertQueryModelQuery.Range = true
	}
}

func Legend(legend string) Option {
	return func(query *Query) {
		query.Builder.Model.NgAlertQueryModelQuery.LegendFormat = legend
	}
}
