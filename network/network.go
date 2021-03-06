package network

import (
	"net"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type NetworkInterface interface {
	EnsureIP(string) error
	RemoveIP(string) error
}

type networkInterface struct {
	card *net.Interface
	link netlink.Link
	arp  *arp
}

func NewNetworkInterfaceFromName(name string) (networkInterface, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return networkInterface{}, errors.Wrapf(err, "fail to open interface %s", name)
	}

	card, err := net.InterfaceByName(name)
	if err != nil {
		return networkInterface{}, errors.Wrapf(err, "fail to find interface %s", name)
	}

	return networkInterface{
		card: card,
		link: link,
		arp:  GetArp(),
	}, nil
}

func (i networkInterface) hasIP(addr *netlink.Addr) (bool, error) {
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

func (i networkInterface) EnsureIP(ip string) error {
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
	err = i.arp.GratuitousArpRequest(GratuitousArpRequest{
		IP:        addr.IP,
		Interface: i.card,
	})
	if err != nil {
		return errors.Wrapf(err, "fail to announce our IP")
	}
	return nil
}

func (i networkInterface) addIP(addr *netlink.Addr) error {
	err := netlink.AddrAdd(i.link, addr)
	if err != nil {
		return errors.Wrap(err, "fail to add address to the interface")
	}

	return nil
}

func (i networkInterface) RemoveIP(ip string) error {
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
