package usage

const (
	// MeasurementId is the unique code of the OpenEBS property in Google Analytics.
	MeasurementId string = "G-TZGP46618W"
	// ApiSecret is the measurement protocol api_secret.
	ApiSecret string = "91JGdTg9QwGn7Y-vvuM4zA"

	// InstallEvent event is sent on pod starts
	InstallEvent string = "install"
	// VolumeProvision event is sent when a volume is created
	VolumeProvision string = "volume-provision"
	// VolumeDeprovision event is sent when a volume is deleted
	VolumeDeprovision string = "volume-deprovision"
	// AppName the application name
	AppName string = "OpenEBS"

	// Replica Event replication
	Replica string = "replica:"
	// DefaultReplicaCount holds the replica count string
	DefaultReplicaCount string = "replica:1"

	// RunningStatus status is running
	RunningStatus string = "running"
	// EventLabelNode holds the string label "nodes"
	EventLabelNode string = "nodes"
	// EventLabelCapacity holds the string label "capacity"
	EventLabelCapacity string = "capacity"
)
