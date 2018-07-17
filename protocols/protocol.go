package protocols

import (
	"gopkg.in/dedis/onet.v2"
	"gopkg.in/dedis/onet.v2/log"
	"github.com/dedis/prifi/prifi-lib/net"
)

// ProtocolName is the name used to register the SDA wrapper protocol with SDA.
const ProtocolName = "DissentProtocol"

//PriFiSDAProtocol is the SDA-protocol struct. It contains the SDA-tree, and a chanel that stops the simulation when it receives a "true"
type DissentProtocol struct {
	*onet.TreeNodeInstance
	configSet     bool
	config        DissentProtocolConfig
	role          DissentRole
	ms            MessageSender
	toHandler     func([]string, []string)
	ResultChannel chan interface{}

	HasStopped       bool //when set to true, the protocol has been stopped by PriFi-lib and should be destroyed
}

//Start is called on the Relay by the service when ChurnHandler decides so
func (p *DissentProtocol) Start() error {

	if !p.configSet {
		log.Fatal("Trying to start PriFi-lib, but config not set !")
	}

	//At the protocol is ready,

	log.Lvl3("Starting Dissent-SDA-Wrapper Protocol")

	//emulate the reception of a ALL_ALL_PARAMETERS with StartNow=true
	msg := new(net.ALL_ALL_PARAMETERS)
	msg.Add("StartNow", true)
	msg.Add("NTrustees", len(p.ms.trustees))
	msg.Add("NClients", len(p.ms.clients))
	msg.Add("PayloadSize", p.config.Toml.PayloadSize)
	msg.Add("DownstreamCellSize", p.config.Toml.CellSizeDown)
	msg.Add("WindowSize", p.config.Toml.RelayWindowSize)
	msg.Add("UseOpenClosedSlots", p.config.Toml.RelayUseOpenClosedSlots)
	msg.Add("UseDummyDataDown", p.config.Toml.RelayUseDummyDataDown)
	msg.Add("ExperimentRoundLimit", p.config.Toml.RelayReportingLimit)
	msg.Add("UseUDP", p.config.Toml.UseUDP)
	msg.Add("DCNetType", p.config.Toml.DCNetType)
	msg.Add("DisruptionProtectionEnabled", p.config.Toml.DisruptionProtectionEnabled)
	msg.Add("OpenClosedSlotsMinDelayBetweenRequests", p.config.Toml.OpenClosedSlotsMinDelayBetweenRequests)
	msg.Add("RelayMaxNumberOfConsecutiveFailedRounds", p.config.Toml.RelayMaxNumberOfConsecutiveFailedRounds)
	msg.Add("RelayProcessingLoopSleepTime", p.config.Toml.RelayProcessingLoopSleepTime)
	msg.Add("RelayRoundTimeOut", p.config.Toml.RelayRoundTimeOut)
	msg.Add("RelayTrusteeCacheLowBound", p.config.Toml.RelayTrusteeCacheLowBound)
	msg.Add("RelayTrusteeCacheHighBound", p.config.Toml.RelayTrusteeCacheHighBound)
	msg.Add("EquivocationProtectionEnabled", p.config.Toml.EquivocationProtectionEnabled)
	msg.ForceParams = true

	p.SendTo(p.TreeNode(), msg)

	return nil
}

// Stop aborts the current execution of the protocol.
func (p *DissentProtocol) Stop() {
	p.HasStopped = true
	p.Shutdown()
}

/**
 * On initialization of the PriFi-SDA-Wrapper protocol, it need to register the PriFi-Lib messages to be able to marshall them.
 * If we forget some messages there, it will crash when PriFi-Lib will call SendToXXX() with this message !
 */
func init() {

	//register the prifi_lib's message with the network lib here
	//network.RegisterMessage(net.ALL_ALL_PARAMETERS{})

	onet.GlobalProtocolRegister(ProtocolName, NewDissentProtocol)
}

// NewPriFiSDAWrapperProtocol creates a bare PrifiSDAWrapper struct.
// SetConfig **MUST** be called on it before it can participate
// to the protocol.
func NewDissentProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	p := &DissentProtocol{
		TreeNodeInstance: n,
		ResultChannel:    make(chan interface{}),
	}

	return p, nil
}

// registerHandlers contains the verbose code
// that registers handlers for all prifi messages.
func (p *DissentProtocol) registerHandlers() error {
	//register handlers
	/*err := p.RegisterHandler(p.Received_ALL_ALL_PARAMETERS_NEW)
	if err != nil {
		return errors.New("couldn't register handler: " + err.Error())
	}*/

	return nil
}
