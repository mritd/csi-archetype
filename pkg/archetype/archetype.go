package archetype

import (
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
)

const (
	driverName = "csi-archetype"
)

type archetype struct {
	name     string
	nodeID   string
	version  string
	endpoint string

	// Add CSI plugin parameters here
	parameter1 string
	parameter2 int
	parameter3 time.Duration

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

func NewCSIDriver(version, nodeID, endpoint, parameter1 string, parameter2 int, parameter3 time.Duration) *archetype {
	logrus.Infof("Driver: %s version: %s", driverName, version)

	// Add some check here
	if parameter1 == "" {
		logrus.Fatal("parameter1 is empty")
	}

	n := &archetype{
		name:     driverName,
		nodeID:   nodeID,
		version:  version,
		endpoint: endpoint,

		parameter1: parameter1,
		parameter2: parameter2,
		parameter3: parameter3,
	}

	// Add access modes for CSI here
	n.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})

	// Add service capabilities for CSI here
	n.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
	})

	return n
}

func (n *archetype) Run() {
	server := NewNonBlockingGRPCServer()
	server.Start(
		n.endpoint,
		NewIdentityServer(n),
		NewControllerServer(n),
		NewNodeServer(n),
	)
	server.Wait()
}

func (n *archetype) AddVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		logrus.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, &csi.VolumeCapability_AccessMode{Mode: c})
	}
	n.cap = vca
}

func (n *archetype) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		logrus.Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}
	n.cscap = csc
}
