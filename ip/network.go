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
	err := m.networkInterface.RemoveIP(m.ip.IP)
	if err != nil {
		log.WithError(err).Error("Fail to de-activate IP")
	}
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
			log.Info("Send gratuitous ARP request")
			err := m.networkInterface.EnsureIP(m.ip.IP)
			if err != nil {
				log.WithError(err).Error("Fail to ensure IP")
			}
		}

		time.Sleep(m.config.ARPGratuitousInterval)
	}

}
