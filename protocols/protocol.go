package protocols

import (
	"gopkg.in/dedis/onet.v2/network"
	"gopkg.in/dedis/onet.v2"
	"gopkg.in/dedis/kyber.v2"
	"github.com/dedis/prifi/prifi-lib/crypto"
	"errors"
)

const ProtocolName = "DissentProtocol"

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

	keyPriv		kyber.Scalar
	keyPub 	kyber.Point

	HasStopped       bool
}

// Stop aborts the current execution of the protocol.
func (p *DissentProtocol) Stop() {
	p.HasStopped = true
	p.Shutdown()
}

func init() {
	network.RegisterMessage(NEW_ROUND{})
	network.RegisterMessage(PUBLIC_KEY{})
	network.RegisterMessage(ALL_ALL_PARAMETERS{})

	onet.GlobalProtocolRegister(ProtocolName, NewDissentProtocol)
}

func (p *DissentProtocol) registerHandlers() error {

	err := p.RegisterHandler(p.Received_NEW_ROUND)
	if err != nil {
		return errors.New("couldn't register handler: " + err.Error())
	}
	err = p.RegisterHandler(p.Received_PUBLIC_KEY)
	if err != nil {
		return errors.New("couldn't register handler: " + err.Error())
	}
	err = p.RegisterHandler(p.Received_ALL_ALL_PARAMETERS)
	if err != nil {
		return errors.New("couldn't register handler: " + err.Error())
	}

	return nil
}

// NewDissentProtocol creates a bare DissentProtocol struct.
// SetConfig **MUST** be called on it before it can participate
// to the protocol.
func NewDissentProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	p := &DissentProtocol{
		TreeNodeInstance: n,
		ResultChannel:    make(chan interface{}),
	}

	return p, nil
}


// SetConfig configures the Dissent node.
// It **MUST** be called in service.newProtocol or before Start().
func (p *DissentProtocol) SetConfigFromDissentService(config *DissentProtocolConfig) {
	p.config = *config
	p.role = config.Role

	ms := p.buildMessageSender(config.Identities)
	p.ms = ms

	p.nClients = len(p.ms.clients)
	p.nTrustees = len(p.ms.trustees)

	p.keyPub, p.keyPriv = crypto.NewKeyPair()

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