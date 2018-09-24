package etcdprobe

import (
	"context"

	etcd "github.com/coreos/etcd/client"
	errgo "gopkg.in/errgo.v1"
)

type EtcdProbe struct {
	name string
	kapi etcd.KeysAPI
}

func NewEtcdProbe(name string, client etcd.KeysAPI) EtcdProbe {
	return EtcdProbe{
		name: name,
		kapi: client,
	}
}

func (p EtcdProbe) Name() string {
	return p.name
}

func (p EtcdProbe) Check() error {
	_, err := p.kapi.Get(context.Background(), "/", nil)
	if err != nil {
		return errgo.Notef(err, "Unable to contact server")
	}
	return nil
}
