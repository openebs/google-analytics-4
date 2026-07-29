## Google Analytics 4
[![CI](https://github.com/openebs/google-analytics-4/actions/workflows/pull_request.yml/badge.svg)](https://github.com/openebs/google-analytics-4/actions/workflows/pull_request.yml)
[![Slack](https://img.shields.io/badge/chat-slack-ff1493.svg?style=flat-square)](https://kubernetes.slack.com/messages/openebs)
[![Community Meetings](https://img.shields.io/badge/Community-Meetings-blue)](https://github.com/openebs/community/blob/HEAD/README.md#community)
[![Go Report](https://goreportcard.com/badge/github.com/openebs/google-analytics-4)](https://goreportcard.com/report/github.com/openebs/google-anaytics-4)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B162%2Fgithub.com%2Fopenebs%2Fgoogle-analytics-4.svg?type=shield&issueType=license)](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fopenebs%2Fgoogle-analytics-4?ref=badge_shield&issueType=license)

Track and monitor your Go programs for free with Google Analytics

The `ga` package is essentially a Go wrapper around the [Google Analytics - Measurement Protocol (Google Analytics 4)](https://developers.google.com/analytics/devguides/collection/protocol/ga4/reference?client_type=gtag)

### Install

```
go get -v github.com/openebs/google-analytics-4
```

### API

Create a new `client` and `Send()` an 'event'.

### Quick Usage

1. Log into GA and create a new property and note its Measurement ID

2. Create a `ga-test.go` file

	``` go
	package main

        import (
            "fmt"

            gaClient "github.com/openebs/google-analytics-4/pkg/client"
            gaEvent "github.com/openebs/google-analytics-4/pkg/event"
        )

        func main() {
            client, err := gaClient.NewMeasurementClient(
                gaClient.WithApiSecret("<api-secret>"),
                gaClient.WithMeasurementId("<measurement-id>"),
                gaClient.WithClientId("<client-id>"),
            )
            if err != nil {
                panic(err)
            }

            event := gaEvent.NewOpenebsEventBuilder().
                Project("OpenEBS").
                K8sVersion("v1.25.15").
                EngineName("test-engine").
                EngineVersion("v1.0.0").
                K8sDefaultNsUid("f5d2a546-19ce-407d-99d4-0655d67e2f76").
                EngineInstaller("helm").
                NodeOs("Ubuntu 20.04.6 LTS").
                NodeArch("linux/amd64").
                NodeKernelVersion("5.4.0-165-generic").
                VolumeName("pvc-b3968e30-9020-4011-943a-7ab338d5f19f").
                VolumeClaimName("openebs-lvmpv").
                Category("volume-deprovision").
                NodeCount("3").
                Build()

            err = client.Send(event)
            if err != nil {
                panic(err)
            }

            fmt.Println("Event fired!")
        }

	```

3. In GA, go to Report > Realtime

4. Run `ga-test.go`

	```
	$ go run ga-test.go
	Event fired!
	```

5. Watch as your event appears

	![foo-ga](https://cloud.githubusercontent.com/assets/633843/5979585/023fc580-a8fd-11e4-803a-956610bcc2e2.png)

### Configuration

The `usage` package reads the following environment variables. All of them are
optional.

| Variable | Description |
| --- | --- |
| `GA_ID` | base64-encoded GA4 Measurement ID for the target property. |
| `GA_KEY` | base64-encoded Measurement Protocol API secret for the target property. |
| `GA_DNS` | DNS server used to resolve `www.google-analytics.com`. When unset, the system resolver is used. |
| `OPENEBS_IO_ANALYTICS_PING_INTERVAL` | Interval between ping events (default `24h`, minimum `1h`). |

`GA_ID` and `GA_KEY` take effect only when **both** are set and decode to a
valid Measurement ID; otherwise the built-in defaults are used and events are
reported to the upstream OpenEBS property.

#### `GA_DNS`

Set this when the pod cannot use the cluster resolver to reach Google Analytics.
The port is optional and defaults to `53`. The host must be an IP address, not a
hostname:

```
GA_DNS=8.8.8.8                       # -> 8.8.8.8:53
GA_DNS=8.8.8.8:5353
GA_DNS=2001:4860:4860::8888          # -> [2001:4860:4860::8888]:53
GA_DNS=[2001:4860:4860::8888]:5353
```

Telemetry is best-effort, so an unparseable `GA_DNS` does not stop the client
from being created: the value is rejected, an error is logged, and the system
resolver is used instead. If events stop arriving after setting `GA_DNS`, check
the logs for `failed to create http client` — a value that is accepted but
points at an unreachable server fails silently at send time.

## License Compliance
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B162%2Fgithub.com%2Fopenebs%2Fgoogle-analytics-4.svg?type=large&issueType=license)](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fopenebs%2Fgoogle-analytics-4?ref=badge_large&issueType=license)
