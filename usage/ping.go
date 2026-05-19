package usage

import (
	"context"
	"fmt"
	"time"

	"github.com/openebs/lib-csi/pkg/common/env"
)

// OpenEBSPingPeriod  ping interval of volume io analytics
var OpenEBSPingPeriod = "OPENEBS_IO_ANALYTICS_PING_INTERVAL"

const (
	// defaultPingPeriod sets the default ping heartbeat interval
	defaultPingPeriod time.Duration = 24 * time.Hour
	// minimumPingPeriod sets the minimum possible configurable
	// heartbeat period, if a value lower than this will be set, the
	// defaultPingPeriod will be used
	minimumPingPeriod time.Duration = 1 * time.Hour
)

// PingCheck sends ping events to Google Analytics on a fixed cadence.
func PingCheck(engineName, category string, pingImmediately bool) {
	PingCheckCtx(context.Background(), engineName, category, pingImmediately)
}

// PingCheckCtx sends ping events to Google Analytics on a fixed cadence,
// returning when ctx is cancelled. If pingImmediately is true, one event is
// sent before the ticker starts; subsequent events fire every GetPingPeriod().
func PingCheckCtx(ctx context.Context, engineName, category string, pingImmediately bool) {
	// Create a new usage field
	u := New()

	if pingImmediately {
		// Ping immediately.
		u.CommonBuild(engineName).
			InstallBuilder(true).
			SetCategory(category).
			Send()
	}

	ticker := time.NewTicker(GetPingPeriod())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Ping periodically.
			u.CommonBuild(engineName).
				InstallBuilder(true).
				SetCategory(category).
				Send()
		}
	}
}

// GetPingPeriod sets the duration of health events, defaults to 24
func GetPingPeriod() time.Duration {
	value := env.GetOrDefault(OpenEBSPingPeriod, fmt.Sprint(defaultPingPeriod))
	duration, _ := time.ParseDuration(value)
	// Sanity checks for setting time duration of health events
	// This way, we are checking for negative and zero time duration and we
	// also have a minimum possible configurable time duration between health events
	if duration < minimumPingPeriod {
		// Avoid corner case when the ENV value is undesirable
		return time.Duration(defaultPingPeriod)
	}

	return duration
}
