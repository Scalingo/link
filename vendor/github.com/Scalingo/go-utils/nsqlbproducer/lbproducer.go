package nsqlbproducer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Scalingo/go-utils/nsqproducer"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	errgo "gopkg.in/errgo.v1"
)

type LBStrategy int

const (
	FallbackStrategy LBStrategy = iota
	RandomStrategy
)

var (
	StrategiesFromName = map[string]LBStrategy{
		"":         FallbackStrategy,
		"fallback": FallbackStrategy,
		"random":   RandomStrategy,
	}
)

// NsqLBProducer a producer that distribute nsq messages across a set of node
// if a node send an error when receiving the message it will try with another node of the set
type NsqLBProducer struct {
	producers []producer
	randInt   func() int
	strategy  LBStrategy
	logger    logrus.FieldLogger
}

type producer struct {
	producer nsqproducer.Producer
	host     Host
}

type Host struct {
	Host string
	Port string
}

func (h Host) String() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

type LBProducerOpts struct {
	Hosts      []Host
	NsqConfig  *nsq.Config
	Logger     logrus.FieldLogger
	SkipLogSet map[string]bool
	Strategy   LBStrategy
}

var _ nsqproducer.Producer = &NsqLBProducer{} // Ensure that NsqLBProducer implements the Producer interface

func New(opts LBProducerOpts) (*NsqLBProducer, error) {
	if len(opts.Hosts) == 0 {
		return nil, fmt.Errorf("A producer must have at least one host")
	}
	lbproducer := &NsqLBProducer{
		producers: make([]producer, len(opts.Hosts)),
		strategy:  opts.Strategy,
	}

	for i, h := range opts.Hosts {
		p, err := nsqproducer.New(nsqproducer.ProducerOpts{
			Host:       h.Host,
			Port:       h.Port,
			NsqConfig:  opts.NsqConfig,
			SkipLogSet: opts.SkipLogSet,
		})

		if err != nil {
			return nil, errors.Wrapf(err, "fail to create producer for host: %s:%s", h.Host, h.Port)
		}

		lbproducer.producers[i] = producer{
			producer: p,
			host:     h,
		}
	}

	switch lbproducer.strategy {
	case FallbackStrategy:
		lbproducer.randInt = alwaysZero
	case RandomStrategy:
		fallthrough
	default:
		lbproducer.randInt = rand.New(rand.NewSource(time.Now().Unix())).Int
	}
	lbproducer.logger = opts.Logger

	return lbproducer, nil
}

func alwaysZero() int {
	return 0
}

func (p *NsqLBProducer) Publish(ctx context.Context, topic string, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randInt() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		producer := p.producers[(i+firstProducer)%len(p.producers)]
		err = producer.producer.Publish(ctx, topic, message)
		if err != nil {
			if p.logger != nil {
				p.logger.WithError(err).WithField("host", producer.host.String()).Error("fail to send nsq message to one nsq node")
			}
		} else {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) DeferredPublish(ctx context.Context, topic string, delay int64, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randInt() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		producer := p.producers[(i+firstProducer)%len(p.producers)]
		err = producer.producer.DeferredPublish(ctx, topic, delay, message)
		if err != nil {
			if p.logger != nil {
				p.logger.WithError(err).WithField("host", producer.host.String()).Error("fail to send nsq message to one nsq node")
			}
		} else {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) Stop() {
	for _, p := range p.producers {
		p.producer.Stop()
	}
}
