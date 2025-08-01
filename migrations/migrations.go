package migrations

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/locker"
	"github.com/Scalingo/link/v3/models"
)

type Migration interface {
	NeedsMigration(ctx context.Context) (bool, error)
	Migrate(ctx context.Context) error
	Name() string
}

type MigrationRunner struct {
	migrations []Migration
}

func NewMigrationRunner(ctx context.Context, cfg config.Config, storage models.Storage, leaseManager locker.EtcdLeaseManager) (MigrationRunner, error) {

	v3tov4, err := NewV3toV4Migration(ctx, cfg, storage)
	if err != nil {
		return MigrationRunner{}, errors.Wrap(ctx, err, "init v3 to v4 migration")
	}
	migrations := []Migration{
		NewV0toV1Migration(cfg.Hostname, leaseManager, storage),
		NewV2toV3Migration(cfg.Hostname, storage),
		v3tov4,
	}

	return MigrationRunner{
		migrations: migrations,
	}, nil
}

func (m MigrationRunner) Run(ctx context.Context) error {
	for _, migration := range m.migrations {
		migrationName := migration.Name()
		ctx, log := logger.WithFieldToCtx(ctx, "migration", migrationName)
		log.Info("Checking if migration is needed")
		needsMigration, err := migration.NeedsMigration(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "check if migration %s is needed", migrationName)
		}
		if !needsMigration {
			log.Info("Migration is not needed")
			continue
		}

		log.Info("Starting migration")
		err = migration.Migrate(ctx)
		if err != nil {
			return errors.Wrapf(ctx, err, "migrate data for %s", migrationName)
		}
		log.Info("Migration completed")
	}

	return nil
}
