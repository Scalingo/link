package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/locker"
	"github.com/looplab/fsm"
)

func (m *manager) setActivated(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: ACTIVATED")
	ip := m.IP()

	if !ip.NoNetwork {
		err := m.networkInterface.EnsureIP(ip.IP)
		if err != nil {
			log.WithError(err).Error("Fail to activate IP")
		}
	}
}

func (m *manager) setStandBy(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: STANDBY")
	ip := m.IP()

	if !ip.NoNetwork {
		err := m.networkInterface.RemoveIP(ip.IP)
		if err != nil {
			log.WithError(err).Error("Fail to de-activate IP")
		}
	}
}

func (m *manager) setFailing(ctx context.Context, _ *fsm.Event) {
	log := logger.Get(ctx)
	log.Info("New state: FAILING")
	ip := m.IP()

	if !ip.NoNetwork {
		err := m.networkInterface.RemoveIP(ip.IP)
		if err != nil {
			log.WithError(err).Error("Fail to de-activate IP")
		}
	}

	err := m.locker.Unlock(ctx)
	if err != nil && err != locker.ErrNotMaster {
		// If we are not master, we can safely ignore this error
		log.WithError(err).Error("Fail to unlock the key")
	}
}

func (m *manager) startArpEnsure(ctx context.Context) {
	var (
		garpCount int
	)
	log := logger.Get(ctx).WithField("process", "arp_ensure")
	for {
		if m.isStopped() {
			return
		}
		currentState := m.Status()
		ip := m.IP()

		if currentState == ACTIVATED && garpCount < m.config.ARPGratuitousCount && !ip.NoNetwork {
			log.Debug("Send gratuitous ARP request")
			err := m.networkInterface.EnsureIP(ip.IP)
			if err != nil {
				log.WithError(err).Error("Fail to ensure IP")
			} else {
				garpCount++
			}
		} else if currentState != ACTIVATED {
			garpCount = 0
		}

		timeToSleep := config.RandomDurationAround(m.config.ARPGratuitousInterval, 0.25)
		time.Sleep(timeToSleep)
	}
}
