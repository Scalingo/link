package network

import (
	"github.com/j-keck/arping"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type NetworkInterface interface {
	AddIP(string) error
	RemoveIP(string) error
}

type networkInterface struct {
	cardName string
	link     netlink.Link
}

func NewNetworkInterfaceFromName(name string) (networkInterface, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return networkInterface{}, errors.Wrapf(err, "fail to open interface %s", name)
	}

	return networkInterface{
		cardName: name,
		link:     link,
	}, nil
}

func (i networkInterface) AddIP(ip string) error {
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return errors.Wrapf(err, "invalid IP: %s", ip)
	}

	err = netlink.AddrAdd(i.link, addr)
	if err != nil {
		return errors.Wrap(err, "fail to add address to the interface")
	}

	err = arping.GratuitousArpOverIfaceByName(addr.IP, i.cardName)
	if err != nil {
		return errors.Wrapf(err, "fail to announce our IP")
	}

	return nil
}

func (i networkInterface) RemoveIP(ip string) error {
	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return errors.Wrapf(err, "invalid IP: %s", ip)
	}

	err = netlink.AddrDel(i.link, addr)
	if err != nil {
		return errors.Wrap(err, "fail to remove address from the interface")
	}

	return nil
}
