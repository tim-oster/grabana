package expr

import "github.com/K-Phoen/sdk"

type Option func(expr *Expr)

type Expr struct {
	Builder          *sdk.NgAlertQuery
	IsAlertCondition bool
}

func New(refId string, options ...Option) *Expr {
	expr := &Expr{
		Builder: &sdk.NgAlertQuery{
			RefId:             refId,
			RelativeTimeRange: sdk.RelativeTimeRange{},
			DatasourceUid:     "__expr__",
			Model: sdk.NgAlertQueryModel{
				NgAlertQueryModelExpression: &sdk.NgAlertQueryModelExpression{
					NgAlertQueryModelExpressionParams: sdk.NgAlertQueryModelExpressionParams{
						Datasource:    sdk.NgAlertQueryModelDatasource{Type: "__expr__", Uid: "__expr__"},
						IntervalMs:    sdk.DefaultIntervalMs,
						MaxDataPoints: sdk.DefaultMaxDataPoints,
						RefId:         refId,
					},
					Cmd: sdk.NgAlertQueryModelCommand{},
				},
			},
		},
		IsAlertCondition: false,
	}
	for _, opt := range options {
		opt(expr)
	}

	if expr.Builder.Model.NgAlertQueryModelExpression.Cmd.Type == "" {
		panic("invalid model command")
	}

	return expr
}

func AlertCondition() Option {
	return func(expr *Expr) {
		expr.IsAlertCondition = true
	}
}

func Math(exprStr string) Option {
	return func(expr *Expr) {
		expr.Builder.Model.NgAlertQueryModelExpression.Cmd = sdk.NgAlertQueryModelCommand{
			Type:        sdk.CommandTypeMath,
			MathCommand: &sdk.MathCommand{Expression: exprStr},
		}
	}
}

func Reduce(refId string, reducer sdk.ReducerFunc, opts ...ReducerOption) Option {
	return func(expr *Expr) {
		cmd := &sdk.ReduceCommand{
			Expression: refId,
			Reducer:    reducer,
		}
		for _, opt := range opts {
			opt(cmd)
		}
		expr.Builder.Model.NgAlertQueryModelExpression.Cmd = sdk.NgAlertQueryModelCommand{
			Type:          sdk.CommandTypeReduce,
			ReduceCommand: cmd,
		}
	}
}

type ReducerOption func(cmd *sdk.ReduceCommand)

func ReduceDropNaN() ReducerOption {
	return func(cmd *sdk.ReduceCommand) {
		cmd.Settings = &sdk.ReduceCommandSettings{
			Mode: sdk.ReduceModeDropNonNumeric,
		}
	}
}

func ReduceReplaceNaN(with float64) ReducerOption {
	return func(cmd *sdk.ReduceCommand) {
		cmd.Settings = &sdk.ReduceCommandSettings{
			Mode:             sdk.ReduceModeReplaceNonNumeric,
			ReplaceWithValue: &with,
		}
	}
}

func Resample(refId, window string, down sdk.ResampleDownSampler, up sdk.ResampleUpSampler) Option {
	return func(expr *Expr) {
		expr.Builder.Model.NgAlertQueryModelExpression.Cmd = sdk.NgAlertQueryModelCommand{
			Type: sdk.CommandTypeResample,
			ResampleCommand: &sdk.ResampleCommand{
				Expression:  refId,
				Window:      window,
				DownSampler: down,
				UpSampler:   up,
			},
		}
	}
}

func Threshold(refId string, opt ThresholdOption) Option {
	return func(expr *Expr) {
		cmd := &sdk.ThresholdCommand{
			Expression: refId,
		}
		opt(cmd)
		expr.Builder.Model.NgAlertQueryModelExpression.Cmd = sdk.NgAlertQueryModelCommand{
			Type:             sdk.CommandTypeThreshold,
			ThresholdCommand: cmd,
		}
	}
}

type ThresholdOption func(cmd *sdk.ThresholdCommand)

func Gt(value float64) ThresholdOption {
	return func(cmd *sdk.ThresholdCommand) {
		cmd.Conditions = append(cmd.Conditions, sdk.ThresholdCondition{
			Evaluator: sdk.ThresholdConditionEval{
				Params: []float64{value},
				Type:   sdk.ThresholdConditionEvalTypeTypeGt,
			},
		})
	}
}

func Lt(value float64) ThresholdOption {
	return func(cmd *sdk.ThresholdCommand) {
		cmd.Conditions = append(cmd.Conditions, sdk.ThresholdCondition{
			Evaluator: sdk.ThresholdConditionEval{
				Params: []float64{value},
				Type:   sdk.ThresholdConditionEvalTypeTypeLt,
			},
		})
	}
}

func WithinRange(from, to float64) ThresholdOption {
	return func(cmd *sdk.ThresholdCommand) {
		cmd.Conditions = append(cmd.Conditions, sdk.ThresholdCondition{
			Evaluator: sdk.ThresholdConditionEval{
				Params: []float64{from, to},
				Type:   sdk.ThresholdConditionEvalTypeTypeWithinRange,
			},
		})
	}
}

func OutsideRange(from, to float64) ThresholdOption {
	return func(cmd *sdk.ThresholdCommand) {
		cmd.Conditions = append(cmd.Conditions, sdk.ThresholdCondition{
			Evaluator: sdk.ThresholdConditionEval{
				Params: []float64{from, to},
				Type:   sdk.ThresholdConditionEvalTypeTypeOutsideRange,
			},
		})
	}
}
