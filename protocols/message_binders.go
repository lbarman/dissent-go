package protocols

import "gopkg.in/dedis/onet.v2/log"

func (p *DissentProtocol) Received_NEW_ROUND(msg NEW_ROUND) error {

	log.Lvl1("Received a New ROUND message")

	return nil
}