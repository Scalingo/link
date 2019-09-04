package network

import (
	"net"
	"sync"

	"github.com/j-keck/arping"
)

/*
 * The library used to make arp request is not thread safe. See issue https://github.com/Scalingo/link/issues/15
 * Here the idea is to have a singleton that ensure only one goroutine will use this library to make ARP requests.
 *
 * With this library if you want to call:
 * err = arping.GratuitousArpOverIfaceByName(addr.IP, i.cardName)
 *
 * Just do:
 * 	err = network.GetARP().GratuitousArpRequest(GratuitousArpRequest{
 * 		IP:        addr.IP,
 * 		Interface: i.cardName,
 * 	})
 */

var arpInstance *arp
var arpOnce sync.Once

type ARP interface {
	GratuitousArp(GratuitousArpRequest) error
}

type GratuitousArpRequest struct {
	IP        net.IP
	Interface *net.Interface
}

type gratuitousArpRequest struct {
	request  GratuitousArpRequest
	response chan error
}

type arp struct {
	requests chan gratuitousArpRequest
}

func GetArp() *arp {
	arpOnce.Do(func() {
		arpInstance = &arp{
			requests: make(chan gratuitousArpRequest),
		}
		go arpInstance.start()
	})

	return arpInstance
}

func (a *arp) GratuitousArpRequest(request GratuitousArpRequest) error {
	errChan := make(chan error)
	a.requests <- gratuitousArpRequest{
		request:  request,
		response: errChan,
	}

	return <-errChan
}

func (a *arp) Stop() {
	close(a.requests)
}

func (a *arp) start() {
	for req := range a.requests {
		err := arping.GratuitousArpOverIface(req.request.IP, *req.request.Interface)
		req.response <- err
	}
}
