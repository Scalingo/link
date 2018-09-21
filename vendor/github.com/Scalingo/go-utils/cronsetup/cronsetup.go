package cronsetup

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-etcd-cron"
	"github.com/Scalingo/go-utils/logger"
	"github.com/coreos/etcd/clientv3"
)

type SetupOpts struct {
	EtcdConfig func() (clientv3.Config, error)
	Jobs       []etcdcron.Job
}

func Setup(ctx context.Context, opts SetupOpts) (func(), error) {
	etcdConfig, err := opts.EtcdConfig()
	if err != nil {
		return nil, fmt.Errorf("fail to get etcd v3 config: %v", err)
	}

	etcdMutexBuilder, err := etcdcron.NewEtcdMutexBuilder(etcdConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to get etcd mutex builder: %v", err)
	}

	funcCtx := func(ctx context.Context, j etcdcron.Job) context.Context {
		log := logger.Get(ctx)
		log = log.WithField("cron-job", j.Name)
		log.Debug("running cron job")
		return logger.ToCtx(ctx, log)
	}

	errorHandler := func(ctx context.Context, j etcdcron.Job, err error) {
		log := logger.Get(ctx)
		log.WithError(err).Error("error when running cron job")
	}

	c, err := etcdcron.New(
		etcdcron.WithEtcdErrorsHandler(errorHandler),
		etcdcron.WithErrorsHandler(errorHandler),
		etcdcron.WithEtcdMutexBuilder(etcdMutexBuilder),
		etcdcron.WithFuncCtx(funcCtx),
	)
	if err != nil {
		return nil, fmt.Errorf("fail to create etcd cron: %v", err)
	}

	for _, job := range opts.Jobs {
		c.AddJob(job)
	}

	log := logger.Get(ctx)
	log.Info("starting etcd-cron")

	c.Start(ctx)
	return func() {
		log.Info("stopping etcd-cron")
		c.Stop()
	}, nil
}
