package network

import (
	"net"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type Interface interface {
	EnsureIP(ip string) error
	RemoveIP(ip string) error
}

type NetInterface struct {
	card *net.Interface
	link netlink.Link
	arp  ARP
}

func NewNetworkInterfaceFromName(name string) (NetInterface, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return NetInterface{}, errors.Wrapf(err, "fail to open interface %s", name)
	}

	card, err := net.InterfaceByName(name)
	if err != nil {
		return NetInterface{}, errors.Wrapf(err, "fail to find interface %s", name)
	}

	return NetInterface{
		card: card,
		link: link,
		arp:  GetArp(),
	}, nil
}

func (i NetInterface) hasIP(addr *netlink.Addr) (bool, error) {
	addrs, err := netlink.AddrList(i.link, netlink.FAMILY_ALL)
	if err != nil {
		return false, errors.Wrap(err, "fail to list interface IPs")
	}

	for _, a := range addrs {
		if a.IP.Equal(addr.IP) {
			return true, nil
		}
	}
	return false, nil
}

func (i NetInterface) EnsureIP(ip string) error {
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return errors.Wrapf(err, "invalid IP: %s", ip)
	}

	has, err := i.hasIP(addr)
	if err != nil {
		return errors.Wrap(err, "fail to check if the IP is present")
	}

	// If the IP has not been configured on the interface
	if !has {
		// Add it
		err := i.addIP(addr)
		if err != nil {
			return errors.Wrap(err, "fail to add IP address")
		}
	}
	// Send the gratuitous ARP request (it wont hurt anyone)
	err = i.arp.GratuitousArp(GratuitousArpRequest{
		IP:        addr.IP,
		Interface: i.card,
	})
	if err != nil {
		return errors.Wrapf(err, "fail to announce our IP")
	}
	return nil
}

func (i NetInterface) addIP(addr *netlink.Addr) error {
	err := netlink.AddrAdd(i.link, addr)
	if err != nil {
		return errors.Wrap(err, "fail to add address to the interface")
	}

	return nil
}

func (i NetInterface) RemoveIP(ip string) error {
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return errors.Wrapf(err, "invalid IP: %s", ip)
	}

	has, err := i.hasIP(addr)
	if err != nil {
		return errors.WithMessage(err, "fail to check if the IP is already present")
	}

	if !has {
		return nil
	}

	err = netlink.AddrDel(i.link, addr)
	if err != nil {
		return errors.Wrap(err, "fail to remove address from the interface")
	}

	return nil
}
