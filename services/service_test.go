package services

import (
	"testing"

	"gopkg.in/dedis/onet.v2/log"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func TestServiceTemplate(t *testing.T) {

	/*
		local := onet.NewLocalTest()
		defer local.CloseAll()
		hosts, roster, _ := local.MakeHELS(5, serviceID)
		log.Lvl1("Roster is", roster)

		var services []onet.Service
		for _, h := range hosts {
			service := local.Services[h.ServerIdentity.ID][serviceID].(onet.Service)
			services = append(services, service)
		}

		services[0].StartTrustee()
		services[1].StartTrustee()
		services[2].StartRelay()
		services[3].StartClient()
		services[4].StartClient()
	*/
}
