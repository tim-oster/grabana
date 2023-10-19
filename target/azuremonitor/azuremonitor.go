package azuremonitor

import "github.com/K-Phoen/sdk"

// Option represents an option that can be used to configure a graphite query.
type Option func(target *AzureMonitor)

// AzureMonitor represents a graphite query.
type AzureMonitor struct {
	Builder *sdk.Target
}

type Aggregation string

const (
	AggregationAvg Aggregation = "Average"
	AggregationMin Aggregation = "Minimum"
	AggregationMax Aggregation = "Maximum"
)

// New creates a new AzureMonitor query.
func New(agg Aggregation, metricNamespace, metricName, region string, options ...Option) *AzureMonitor {
	target := &AzureMonitor{
		Builder: &sdk.Target{
			QueryType: "Azure Monitor",
			AzureMonitor: &sdk.AzureMonitorTarget{
				Aggregation:     string(agg),
				MetricName:      metricName,
				MetricNamespace: metricNamespace,
				Region:          region,
				Resources:       nil,
				TimeGrain:       "auto",
			},
		},
	}

	for _, opt := range options {
		opt(target)
	}

	return target
}

func Resource(subscription, resourceGroup, resourceName string) Option {
	return func(target *AzureMonitor) {
		am := target.Builder.AzureMonitor
		am.Resources = append(am.Resources, sdk.AzureMonitorTargetResource{
			MetricNamespace: am.MetricNamespace,
			Region:          am.Region,
			ResourceGroup:   resourceGroup,
			ResourceName:    resourceName,
			Subscription:    subscription,
		})
	}
}

func TimeGain(timeGrain string) Option {
	return func(target *AzureMonitor) {
		target.Builder.AzureMonitor.TimeGrain = timeGrain
	}
}
