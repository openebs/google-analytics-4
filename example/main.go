package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	gaClient "github.com/openebs/google-analytics-4/pkg/client"
	gaEvent "github.com/openebs/google-analytics-4/pkg/event"
)

func httpClientWithDns(dns string) (*http.Client, error) {
	if _, _, err := net.SplitHostPort(dns); err != nil {
		return nil, fmt.Errorf("invalid dns address %q: must be host:port", dns)
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, dns)
			},
		},
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = dialer.DialContext
	return &http.Client{Transport: tr}, nil
}

func main() {

	opts := []gaClient.MeasurementClientOption{
		gaClient.WithApiSecret("<api-secret>"),
		gaClient.WithMeasurementId("<measurement-id>"),
		gaClient.WithClientId("1b803d56-fde0-4f1e-ab64-ccb22509ae9f"),
	}

	// Only configure a custom HTTP client (with a DNS override) when a DNS
	// address is set. When it's empty, skip the option so the measurement
	// client falls back to its default HTTP transport.
	dns := "<dns>"
	if dns != "" {
		httpClient, err := httpClientWithDns(dns)
		if err != nil {
			panic(err)
		}
		opts = append(opts, gaClient.WithHttpClient(httpClient))
	}

	client, err := gaClient.NewMeasurementClient(opts...)
	if err != nil {
		panic(err)
	}

	event := gaEvent.NewOpenebsEventBuilder().
		Project("OpenEBS").
		K8sVersion("v1.25.15").
		EngineName("test-engine").
		EngineVersion("v1.0.0").
		K8sDefaultNsUid("1b803d56-fde0-4f1e-ab64-ccb22509ae9f").
		EngineInstaller("helm").
		NodeOs("Ubuntu 20.04.6 LTS").
		NodeArch("linux/amd64").
		NodeKernelVersion("5.4.0-165-generic").
		VolumeName("pvc-b3968e30-9020-4011-943a-7ab338d5f19f").
		VolumeClaimName("openebs-lvmpv").
		Category("volume_deprovision").
		NodeCount("2").
		VolumeCapacity("19238457924875977657").
		ReplicaCount("1000").
		Build()

	err = client.Send(event)
	if err != nil {
		panic(err)
	}

	fmt.Println("Event fired!")
}
