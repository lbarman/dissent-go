package services

import (
	"fmt"
	"io/ioutil"
	"os"

	dissent_protocol "github.com/lbarman/dissent-go/protocols"

	"gopkg.in/dedis/onet.v2/app"
	"gopkg.in/dedis/onet.v2/log"
	"gopkg.in/dedis/onet.v2/network"
)


//Set the config, from the dissent.toml. Is called by sda/app.
func (s *ServiceState) SetConfigFromToml(config *dissent_protocol.DissentTomlConfig) {
	log.Lvl3("Setting Dissent configuration...")
	log.Lvlf3("%+v\n", config)
	s.prifiTomlConfig = config
}

// tryLoad tries to load the configuration and updates if a configuration
// is found, else it returns an error.
func (s *ServiceState) tryLoad() error {
	configFile := s.path + "/identity.bin"
	b, err := ioutil.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error while reading %s: %s", configFile, err)
	}
	if len(b) > 0 {
		_, msg, err := network.Unmarshal(b, s.Suite())
		if err != nil {
			return fmt.Errorf("Couldn't unmarshal: %s", err)
		}
		log.Lvl3("Successfully loaded")
		s.Storage = msg.(*Storage)
	}
	return nil
}

// mapIdentities reads the group configuration to assign PriFi roles
// to server addresses and returns them with the server
// identity of the relay.
func mapIdentities(group *app.Group) (*network.ServerIdentity, []*network.ServerIdentity) {
	trustees := make([]*network.ServerIdentity, 0)
	var relay *network.ServerIdentity

	// Read the description of the nodes in the config file to assign them PriFi roles.
	nodeList := group.Roster.List
	for i := 0; i < len(nodeList); i++ {
		si := nodeList[i]
		nodeDescription := group.GetDescription(si)

		if nodeDescription == "relay" {
			relay = si
		} else if nodeDescription == "trustee" {
			trustees = append(trustees, si)
		}
	}

	return relay, trustees
}
