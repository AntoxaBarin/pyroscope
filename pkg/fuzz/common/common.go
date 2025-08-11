package common

import (
	"flag"
	"github.com/grafana/pyroscope/api/gen/proto/go/fuzz"
	"github.com/grafana/pyroscope/pkg/distributor/ingestlimits"
	"github.com/grafana/pyroscope/pkg/validation"
	"time"
)

func Overrides(message *fuzz.Limits) *validation.Overrides {
	var limits = validation.Limits{}
	if message == nil {
		limits.RegisterFlags(flag.NewFlagSet("", flag.ContinueOnError))
	} else {
		limits = validation.Limits{
			IngestionRateMB:                  float64(message.IngestionRateMB),
			IngestionBurstSizeMB:             float64(message.IngestionBurstSizeMB),
			MaxLabelNameLength:               int(message.MaxLabelNameLength),
			MaxLabelValueLength:              int(message.MaxLabelValueLength),
			MaxLabelNamesPerSeries:           int(message.MaxLabelNamesPerSeries),
			MaxSessionsPerSeries:             int(message.MaxSessionsPerSeries),
			EnforceLabelsOrder:               message.EnforceLabelsOrder,
			MaxProfileSizeBytes:              int(message.MaxProfileSizeBytes),
			MaxProfileStacktraceSamples:      int(message.MaxProfileStacktraceSamples),
			MaxProfileStacktraceSampleLabels: int(message.MaxProfileStacktraceSampleLabels),
			MaxProfileStacktraceDepth:        int(message.MaxProfileStacktraceDepth),
			MaxProfileSymbolValueLength:      int(message.MaxProfileSymbolValueLength),

			//IngestionLimit:                   message.Limits.IngestionLimit,
			//DistributorSampling:              message.Limits.DistributorSampling,
		}
		if message.IngestionLimit != nil {
			limits.IngestionLimit = &ingestlimits.Config{
				PeriodType:     message.IngestionLimit.PeriodType,
				PeriodLimitMb:  int(message.IngestionLimit.PeriodLimitMb),
				LimitResetTime: message.IngestionLimit.LimitResetTime,
				LimitReached:   message.IngestionLimit.LimitReached,
				Sampling:       ingestlimits.SamplingConfig{},
				UsageGroups:    nil,
			}
			if message.IngestionLimit.Sampling != nil {
				limits.IngestionLimit.Sampling.NumRequests = int(message.IngestionLimit.Sampling.NumRequests)
				//todo this probably will not work without time mocking
				limits.IngestionLimit.Sampling.Period = time.Duration(int(time.Millisecond) * int(message.IngestionLimit.Sampling.PeriodMillis))
			}
			limits.IngestionLimit.UsageGroups = make(map[string]ingestlimits.UsageGroup)
			for k, v := range message.IngestionLimit.UsageGroups {
				if v == nil {
					continue
				}
				limits.IngestionLimit.UsageGroups[k] = ingestlimits.UsageGroup{
					PeriodLimitMb: int(v.PeriodLimitMb), LimitReached: v.LimitReached,
				}
			}
		}
	}
	//IngestionLimit(tenantID string) *ingestlimits.Config
	//DistributorSampling(tenantID string) *sampling.Config
	//IngestionTenantShardSize(tenantID string) int

	//IngestionRelabelingRules(tenantID string) []*relabel.Config
	//DistributorUsageGroups(tenantID string) *validation.UsageGroupConfig

	res, err := validation.NewOverrides(limits, nil)
	if err != nil {
		panic(err)
	}
	return res
}
