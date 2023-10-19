package azure

import (
	"github.com/K-Phoen/grabana/target/azuremonitor"
	"github.com/K-Phoen/sdk"
)

type Query struct {
	Builder          *sdk.NgAlertQuery
	IsAlertCondition bool
}

func New(refId string, target *sdk.AzureMonitorTarget) *Query {
	var sub string
	for _, r := range target.Resources {
		if sub != "" && r.Subscription != sub {
			panic("targets use different subscriptions")
		}
		sub = r.Subscription
	}
	if sub == "" {
		panic("no subscription in targets or no targets")
	}

	return &Query{
		Builder: &sdk.NgAlertQuery{
			RefId: refId,
			RelativeTimeRange: sdk.RelativeTimeRange{
				FromSeconds: 3600,
				ToSeconds:   0,
			},
			QueryType: "Azure Monitor",
			Model: sdk.NgAlertQueryModel{
				NgAlertQueryModelCustom: &queryModel{
					RefId:        refId,
					AzureMonitor: target,
					QueryType:    "Azure Monitor",
					Subscription: sub,
				},
			},
		},
		IsAlertCondition: false,
	}
}

func MakeTarget(agg azuremonitor.Aggregation, metricNamespace, metricName, region string, options ...azuremonitor.Option) *sdk.AzureMonitorTarget {
	return azuremonitor.New(agg, metricNamespace, metricName, region, options...).Builder.AzureMonitor
}

type queryModel struct {
	RefId        string                  `json:"refId"`
	AzureMonitor *sdk.AzureMonitorTarget `json:"azureMonitor"`
	QueryType    string                  `json:"queryType"`
	Subscription string                  `json:"subscription"`
}
