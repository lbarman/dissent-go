package protocols

import (
	"errors"
	"strconv"

	"github.com/dedis/prifi/prifi-lib/net"
	"gopkg.in/dedis/onet.v2"
	"gopkg.in/dedis/onet.v2/log"
)

//MessageSender is the struct we need to give PriFi-Lib so it can send messages.
//It needs to implement the "MessageSender interface" defined in prifi_lib/prifi.go
type MessageSender struct {
	tree       *onet.TreeNodeInstance
	client0    *onet.TreeNode
	clients    map[int]*onet.TreeNode
	trustees   map[int]*onet.TreeNode
}

// buildMessageSender creates a MessageSender struct
// given a mep between server identities and PriFi identities.
func (p *DissentProtocol) buildMessageSender(identities map[string]DissentIdentity) MessageSender {
	nodes := p.List() // Has type []*onet.TreeNode
	trustees := make(map[int]*onet.TreeNode)
	clients := make(map[int]*onet.TreeNode)
	trusteeID := 0
	clientID := 0
	var relay *onet.TreeNode

	for i := 0; i < len(nodes); i++ {
		identifier := nodes[i].ServerIdentity.Public.String()
		id, ok := identities[identifier]
		port, _ := strconv.Atoi(nodes[i].ServerIdentity.Address.Port())
		portForFastChannel := port + 3

		log.Lvl3("Found identity", identifier, " -> ", port, portForFastChannel)

		if !ok {
			log.Lvl3("Skipping unknow node with address", identifier)
			continue
		}
		switch id.Role {
		case Client:
			clients[clientID] = nodes[i] //TODO : wrong
			clientID++
		case Trustee:
			trustees[trusteeID] = nodes[i]
			trusteeID++
		case Client0:
			if relay == nil {
				relay = nodes[i]
			} else {
				log.Fatal("Multiple relays")
			}
		}
	}

	return MessageSender{p.TreeNodeInstance, relay, clients, trustees}
}

//SendToClient sends a message to client i, or fails if it is unknown
func (ms MessageSender) FastSendToClient(i int, msg *net.REL_CLI_DOWNSTREAM_DATA) error {

	if client, ok := ms.clients[i]; ok {
		log.Lvl5("Sending a message to client ", i, " (", client.Name(), ") - ", msg)
		return ms.tree.SendTo(client, msg)
	}

	e := "Client " + strconv.Itoa(i) + " is unknown !"
	log.Error(e)
	return errors.New(e)
}

//SendToClient sends a message to client i, or fails if it is unknown
func (ms MessageSender) SendToClient(i int, msg interface{}) error {

	if client, ok := ms.clients[i]; ok {
		log.Lvl5("Sending a message to client ", i, " (", client.Name(), ") - ", msg)
		return ms.tree.SendTo(client, msg)
	}

	e := "Client " + strconv.Itoa(i) + " is unknown !"
	log.Error(e)
	return errors.New(e)
}

//SendToTrustee sends a message to trustee i, or fails if it is unknown
func (ms MessageSender) SendToTrustee(i int, msg interface{}) error {

	if trustee, ok := ms.trustees[i]; ok {
		log.Lvl5("Sending a message to trustee ", i, " (", trustee.Name(), ") - ", msg)
		return ms.tree.SendTo(trustee, msg)
	}

	e := "Trustee " + strconv.Itoa(i) + " is unknown !"
	log.Error(e)
	return errors.New(e)
}
