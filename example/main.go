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
		gaClient.WithClientId("1b803d56-fde0-4f1e-ab64-ccb22509ae9f"),
	)
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
		Action("replica").
		Label("Capacity").
		Value("2").
		Build()

	err = client.Send(event)
	if err != nil {
		panic(err)
	}

	fmt.Println("Event fired!")
}
