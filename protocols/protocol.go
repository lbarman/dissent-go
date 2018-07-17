package protocols

import (
	"gopkg.in/dedis/onet.v2"
	"gopkg.in/dedis/onet.v2/log"
	"errors"
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

	nClients int
	nTrustees int

	HasStopped       bool //when set to true, the protocol has been stopped by PriFi-lib and should be destroyed
}

//Start is called on the Client0 by the service when ChurnHandler decides so
func (p *DissentProtocol) Start() error {

	if !p.configSet {
		log.Fatal("Trying to start Dissent protocol, but config not set !")
	}

	//At the protocol is ready,

	log.Lvl3("Starting Dissent protocol (", len(p.ms.clients), "c", len(p.ms.trustees), "t)")

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
	err := p.RegisterHandler(p.Received_NEW_ROUND)
	if err != nil {
		return errors.New("couldn't register handler: " + err.Error())
	}

	return nil
}
