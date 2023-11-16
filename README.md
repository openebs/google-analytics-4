## Google Analytics 4

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
                Action("replica:2").
                Label("Capacity").
                Value("2Gi").
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