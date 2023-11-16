/*
Copyright 2023 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package event

type OpenebsEventOption func(*OpenebsEvent) error

// OpenebsEvent Hit Type
type OpenebsEvent struct {
	// Specify Project name, ex OpenEBS
	Project string `json:"project"`
	// K8s Version, ex v1.25.15
	K8sVersion string `json:"k8s_version"`
	// Name of the engine, ex lvm-localpv
	EngineName string `json:"engine_name"`
	// Version of the engine, ex lvm-v1.3.0-e927123:11-08-2023-e927123
	EngineVersion string `json:"engine_version"`
	// Uid of the default k8s ns, ex f5d2a546-19ce-407d-99d4-0655d67e2f76
	K8sDefaultNsUid string `json:"k8s_default_ns_uid"`
	// Installer of the app, ex lvm-localpv-helm
	EngineInstaller string `json:"engine_installer"`
	// Machine's os, ex Ubuntu 20.04.6 LTS
	NodeOs string `json:"node_os"`
	// Machine's kernel version, ex 5.4.0-165-generic
	NodeKernelVersion string `json:"node_kernel_version"`
	// Machine's arch, ex linux/amd64
	NodeArch string `json:"node_arch"`
	// Name of the pv object, example `pvc-b3968e30-9020-4011-943a-7ab338d5f19f`
	VolumeName string `json:"vol_name"`
	// Name of the pvc object, example `openebs-lvmpv`
	VolumeClaimName string `json:"vol_claim_name"`
	// Category of event, i.e install, volume-provision
	Category string `json:"event_category"`
	// Action of the event, i.e running, replica:1
	Action string `json:"event_action"`
	// Label for the event, i.e nodes, capacity
	Label string `json:"event_label"`
	// Value for the label, i.e 4, 2
	Value string `json:"event_value"`
}

// OpenebsEventBuilder is builder for OpenebsEvent
type OpenebsEventBuilder struct {
	openebsEvent *OpenebsEvent
}

func NewOpenebsEventBuilder() *OpenebsEventBuilder {
	openebsEvent := &OpenebsEvent{}
	b := &OpenebsEventBuilder{openebsEvent: openebsEvent}
	return b
}

func (b *OpenebsEventBuilder) Project(project string) *OpenebsEventBuilder {
	b.openebsEvent.Project = project
	return b
}

func (b *OpenebsEventBuilder) K8sVersion(k8sVersion string) *OpenebsEventBuilder {
	b.openebsEvent.K8sVersion = k8sVersion
	return b
}

func (b *OpenebsEventBuilder) EngineName(engineName string) *OpenebsEventBuilder {
	b.openebsEvent.EngineName = engineName
	return b
}

func (b *OpenebsEventBuilder) EngineVersion(engineVersion string) *OpenebsEventBuilder {
	b.openebsEvent.EngineVersion = engineVersion
	return b
}

func (b *OpenebsEventBuilder) K8sDefaultNsUid(k8sDefaultNsUid string) *OpenebsEventBuilder {
	b.openebsEvent.K8sDefaultNsUid = k8sDefaultNsUid
	return b
}

func (b *OpenebsEventBuilder) EngineInstaller(engineInstaller string) *OpenebsEventBuilder {
	b.openebsEvent.EngineInstaller = engineInstaller
	return b
}

func (b *OpenebsEventBuilder) NodeOs(nodeOs string) *OpenebsEventBuilder {
	b.openebsEvent.NodeOs = nodeOs
	return b
}

func (b *OpenebsEventBuilder) NodeKernelVersion(nodeKernelVersion string) *OpenebsEventBuilder {
	b.openebsEvent.NodeKernelVersion = nodeKernelVersion
	return b
}

func (b *OpenebsEventBuilder) NodeArch(nodeArch string) *OpenebsEventBuilder {
	b.openebsEvent.NodeArch = nodeArch
	return b
}

func (b *OpenebsEventBuilder) VolumeName(volumeName string) *OpenebsEventBuilder {
	b.openebsEvent.VolumeName = volumeName
	return b
}

func (b *OpenebsEventBuilder) VolumeClaimName(volumeClaimName string) *OpenebsEventBuilder {
	b.openebsEvent.VolumeClaimName = volumeClaimName
	return b
}

func (b *OpenebsEventBuilder) Category(category string) *OpenebsEventBuilder {
	b.openebsEvent.Category = category
	return b
}

func (b *OpenebsEventBuilder) Action(action string) *OpenebsEventBuilder {
	b.openebsEvent.Action = action
	return b
}

func (b *OpenebsEventBuilder) Label(label string) *OpenebsEventBuilder {
	b.openebsEvent.Label = label
	return b
}

func (b *OpenebsEventBuilder) Value(value string) *OpenebsEventBuilder {
	b.openebsEvent.Value = value
	return b
}

func (b *OpenebsEventBuilder) Build() *OpenebsEvent {
	return b.openebsEvent
}
