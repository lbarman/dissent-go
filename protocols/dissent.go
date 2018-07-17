package protocols

import (
	prifi_lib "github.com/dedis/prifi/prifi-lib"
	"gopkg.in/dedis/onet.v2/log"
	"gopkg.in/dedis/onet.v2/network"
)

//DissentRole is the type of the enum to qualify the role of a SDA node (Relay, Client, Trustee)
type DissentRole int

//The possible states of a SDA node, of type DissentRole
const (
	Relay   DissentRole = iota
	Client
	Trustee
)

//PriFiIdentity is the identity (role + ID)
type PriFiIdentity struct {
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
	Identities            map[string]PriFiIdentity
	Role                  DissentRole
}

// SetConfig configures the PriFi node.
// It **MUST** be called in service.newProtocol or before Start().
func (p *PriFiSDAProtocol) SetConfigFromDissentService(config *DissentProtocolConfig) {
	p.config = *config
	p.role = config.Role

	ms := p.buildMessageSender(config.Identities)
	p.ms = ms

	//sanity check
	switch config.Role {
	case Trustee:
		if ms.relay == nil {
			log.Fatal("Relay is not reachable (I'm a trustee, and I need it) !")
		}
	case Client:
		if ms.relay == nil {
			log.Fatal("Relay is not reachable (I'm a client, and I need it) !")
		}
	case Relay:
		if len(ms.clients) < 1 {
			log.Fatal("Less than one client reachable (I'm a relay, and there's no use starting the protocol) !")
		}
		if len(ms.trustees) < 1 {
			log.Fatal("No trustee reachable (I'm a relay, and I cannot start the protocol) !")
		}
	}

	experimentResultChan := p.ResultChannel

	switch config.Role {
	case Relay:
		relayOutputEnabled := config.Toml.RelayDataOutputEnabled
		p.prifiLibInstance = prifi_lib.NewPriFiRelay(relayOutputEnabled,
			config.RelaySideSocksConfig.DownstreamChannel,
			config.RelaySideSocksConfig.UpstreamChannel,
			experimentResultChan,
			p.handleTimeout,
			ms)
	case Trustee:
		p.prifiLibInstance = prifi_lib.NewPriFiTrustee(config.Toml.TrusteeNeverSlowDown,
			config.Toml.TrusteeAlwaysSlowDown,
			config.Toml.TrusteeSleepTimeBetweenMessages,
			ms)

	case Client:
		doLatencyTests := config.Toml.DoLatencyTests
		clientDataOutputEnabled := config.Toml.ClientDataOutputEnabled
		p.prifiLibInstance = prifi_lib.NewPriFiClient(doLatencyTests,
			clientDataOutputEnabled,
			config.ClientSideSocksConfig.UpstreamChannel,
			config.ClientSideSocksConfig.DownstreamChannel,
			config.Toml.ReplayPCAP,
			config.Toml.PCAPFolder,
			ms)
	}

	p.registerHandlers()

	p.configSet = true
}

// SetTimeoutHandler sets the function that will be called on round timeout
// if the protocol runs as the relay.
func (p *PriFiSDAProtocol) SetTimeoutHandler(handler func([]string, []string)) {
	p.toHandler = handler
}
