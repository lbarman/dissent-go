package protocols

import (
	"gopkg.in/dedis/onet.v2/network"
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

// SetConfig configures the PriFi node.
// It **MUST** be called in service.newProtocol or before Start().
func (p *DissentProtocol) SetConfigFromDissentService(config *DissentProtocolConfig) {
	p.config = *config
	p.role = config.Role

	ms := p.buildMessageSender(config.Identities)
	p.ms = ms

	switch config.Role {
	case Client0:
		/*relayOutputEnabled := config.Toml.RelayDataOutputEnabled
		p.prifiLibInstance = prifi_lib.NewPriFiRelay(relayOutputEnabled,
			config.RelaySideSocksConfig.DownstreamChannel,
			config.RelaySideSocksConfig.UpstreamChannel,
			experimentResultChan,
			p.handleTimeout,
			ms)*/
	case Trustee:

	case Client:
	}

	p.registerHandlers()

	p.configSet = true
}