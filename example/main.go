package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/openebs/google-analytics-4/internal/dnsaddr"
	gaClient "github.com/openebs/google-analytics-4/pkg/client"
	gaEvent "github.com/openebs/google-analytics-4/pkg/event"
)

// httpClientWithDns returns an HTTP client whose resolver queries the given DNS
// server instead of the system default. See dnsaddr.Normalize for the accepted
// address formats.
func httpClientWithDns(dns string) (*http.Client, error) {
	addr, err := dnsaddr.Normalize(dns)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, network, addr)
			},
		},
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.DialContext = dialer.DialContext
	return &http.Client{Transport: tr}, nil
}

func main() {

	opts := []gaClient.MeasurementClientOption{
		// pkg/client takes these verbatim -- do not base64-encode them. Only the
		// usage package's GA_KEY/GA_ID env vars are base64-decoded.
		gaClient.WithApiSecret("<api-secret>"),
		gaClient.WithMeasurementId("<measurement-id>"),
		gaClient.WithClientId("1b803d56-fde0-4f1e-ab64-ccb22509ae9f"),
	}

	// Only configure a custom HTTP client (with a DNS override) when a DNS
	// address is set. When it's empty -- the default when GA_DNS is unset --
	// skip the option so the measurement client falls back to its default HTTP
	// transport. Export GA_DNS=8.8.8.8 (the port is optional) to resolve
	// through a specific server instead.
	dns := strings.TrimSpace(os.Getenv("GA_DNS"))
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
