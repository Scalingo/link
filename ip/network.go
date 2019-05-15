package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/looplab/fsm"
)

func (m *manager) setActivated(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: ACTIVATED")
	err := m.networkInterface.EnsureIP(m.ip.IP)
	if err != nil {
		log.WithError(err).Error("Fail to activate IP")
	}
}

func (m *manager) setStandBy(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: STANDBY")
	// Here we don't want to do anything.
	// - if we came from a success state, we should keep the IP but not advertise it
	//   we keep the IP to prevent broken connections. We still want to respond to this
	//   IP even if another host is master.
	// - If we came from failing state connections were already broken, no need to add the IP back.
}

func (m *manager) setFailing(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: FAILING")

	err := m.networkInterface.RemoveIP(m.ip.IP)
	if err != nil {
		log.WithError(err).Error("Fail to de-activate IP")
	}

	err = m.locker.Stop(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to stop locker")
	}
}

func (m *manager) startArpEnsure(ctx context.Context) {
	log := logger.Get(ctx).WithField("process", "arp_ensure")
	for {
		if m.isStopping() {
			return
		}

		if m.stateMachine.Current() == ACTIVATED {
			log.Debug("Send gratuitous ARP request")
			err := m.networkInterface.EnsureIP(m.ip.IP)
			if err != nil {
				log.WithError(err).Error("Fail to ensure IP")
			}
		}

		time.Sleep(m.config.ARPGratuitousInterval)
	}
}
