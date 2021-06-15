package watcher

import (
	"context"
	"sync"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"go.etcd.io/etcd/v3/clientv3"
)

type KeyChangedCallback func(ctx context.Context)

type EtcdWatcher struct {
	client      clientv3.Watcher
	prefix      string
	cancelLock  *sync.Mutex
	cancelWatch context.CancelFunc
	callback    KeyChangedCallback
}

type Watcher interface {
	Start(context.Context) error
	Stop(context.Context) error
}

func NewWatcher(client clientv3.Watcher, prefix string, callback KeyChangedCallback) Watcher {
	return &EtcdWatcher{
		client:     client,
		prefix:     prefix,
		cancelLock: &sync.Mutex{},
		callback:   callback,
	}
}

func (w *EtcdWatcher) Start(ctx context.Context) error {
	go w.mainEtcdLoop(ctx)
	return nil
}

func (w *EtcdWatcher) Stop(ctx context.Context) error {
	if w.cancelWatch != nil {
		w.cancelWatch()
	}
	return nil
}

// worker does all the heavy lifting on the etcd side. It returns true if it finished successfully or false if there was an error (and should be retried)
func (w *EtcdWatcher) readEtcdWatchEvents(ctx context.Context, resp clientv3.WatchChan) bool {
	log := logger.Get(ctx).WithField("prefix", w.prefix)
	for change := range resp {
		if change.Err() != nil {
			log.WithError(change.Err()).Error("fail to subscribe to changes")
			return false
		}
		if len(change.Events) == 0 {
			continue
		}

		log.Info("Network listeners changed")
		if w.callback != nil {
			go w.callback(ctx)
		}
	}
	log.Infof("Watcher for %s stopped", w.prefix)
	return true
}

func (w *EtcdWatcher) mainEtcdLoop(ctx context.Context) {
	log := logger.Get(ctx)
	for {
		ctx, cancel := context.WithCancel(ctx)
		w.cancelLock.Lock()
		w.cancelWatch = cancel
		w.cancelLock.Unlock()
		respChan := w.client.Watch(ctx, w.prefix, clientv3.WithPrefix())
		stopped := w.readEtcdWatchEvents(ctx, respChan)
		if stopped {
			return
		}

		log.Info("Connection to ETCD lost, waiting 1s and retrying")
		time.Sleep(1 * time.Second)
	}
}
