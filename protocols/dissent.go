package protocols

import (
	"gopkg.in/dedis/onet.v2/network"
	"gopkg.in/dedis/onet.v2/log"
)

//DissentRole is the type of the enum to qualify the role of a SDA node (Client0, Client, Trustee)
type DissentRole int

//The possible states of a SDA node, of type DissentRole
const (
	Client0 DissentRole = iota
	Client
	Trustee
)

//PriFiIdentity is the identity (role + ID)
type DissentIdentity struct {
	Role     DissentRole
	ID       int
	ServerID *network.ServerIdentity
}

//The configuration read in prifi.toml
type DissentTomlConfig struct {
	EnforceSameVersionOnNodes               bool
	ForceConsoleColor                       bool
	OverrideLogLevel                        int
	ClientDataOutputEnabled                 bool
	RelayDataOutputEnabled                  bool
	PayloadSize                             int
	CellSizeDown                            int
	RelayWindowSize                         int
	RelayUseOpenClosedSlots                 bool
	RelayUseDummyDataDown                   bool
	RelayReportingLimit                     int
	UseUDP                                  bool
	DoLatencyTests                          bool
	SocksServerPort                         int
	SocksClientPort                         int
	ProtocolVersion                         string
	DCNetType                               string
	ReplayPCAP                              bool
	PCAPFolder                              string
	TrusteeSleepTimeBetweenMessages         int
	TrusteeAlwaysSlowDown                   bool
	TrusteeNeverSlowDown                    bool
	SimulDelayBetweenClients                int
	DisruptionProtectionEnabled             bool
	EquivocationProtectionEnabled           bool // not linked in the back
	OpenClosedSlotsMinDelayBetweenRequests  int
	RelayMaxNumberOfConsecutiveFailedRounds int
	RelayProcessingLoopSleepTime            int
	RelayRoundTimeOut                       int
	RelayTrusteeCacheLowBound               int
	RelayTrusteeCacheHighBound              int
	VerboseIngressEgressServers             bool
}

//PriFiSDAWrapperConfig is all the information the SDA-Protocols needs. It contains the network map of identities, our role, and the socks parameters if we are the corresponding role
type DissentProtocolConfig struct {
	Toml                  *DissentTomlConfig
	Identities            map[string]DissentIdentity
	Role                  DissentRole
}


//Start is called on the Client0 by the service when ChurnHandler decides so
func (p *DissentProtocol) Start() error {

	if !p.configSet {
		log.Fatal("Trying to start Dissent protocol, but config not set !")
	}
	log.Lvl3("Starting Dissent protocol (", p.nClients, "clients &", p.nTrustees, "trustees)")

	//broadcast the parameters
	message := &ALL_ALL_PARAMETERS{ NClients: p.nClients, NTrustees:p.nTrustees}

	i := 0
	for i < p.nClients {
		p.ms.SendToClient(i, message)
		i++
	}
	i = 0
	for i < p.nTrustees {
		p.ms.SendToTrustee(i, message)
		i++
	}

	return nil
}

func (p *DissentProtocol) Received_ALL_ALL_PARAMETERS(msg Struct_ALL_ALL_PARAMETERS) error {

	log.Lvl1("Received_ALL_ALL_PARAMETERS", p.nClients, p.nTrustees)

	p.nClients = msg.NClients
	p.nTrustees = msg.NTrustees

	//broadcast my key
	message := &PUBLIC_KEY{ Key: p.keyPub}

	i := 0
	for i < p.nClients {
		p.ms.SendToClient(i, message)
		i++
	}
	i = 0
	for i < p.nTrustees {
		p.ms.SendToTrustee(i, message)
		i++
	}

	return nil
}

func (p *DissentProtocol) Received_PUBLIC_KEY(msg Struct_PUBLIC_KEY) error {

	log.Lvl1("Received_PUBLIC_KEY from", msg.ServerIdentity)

	return nil
}

func (p *DissentProtocol) Received_NEW_ROUND(msg Struct_NEW_ROUND) error {

	log.Lvl1("Received_NEW_ROUND")

	return nil
}