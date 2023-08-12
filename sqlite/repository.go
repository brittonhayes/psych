package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"log/slog"

	"github.com/brittonhayes/therapy"
	"github.com/brittonhayes/therapy/sqlite/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/migrate"
)

type repository struct {
	logger *slog.Logger
	m      *migrate.Migrator
	db     *bun.DB
}

func NewRepository(connection string, logger *slog.Logger) therapy.Repository {
	sqldb, err := sql.Open(sqliteshim.ShimName, connection)
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	db.RegisterModel((*therapy.Therapist)(nil))

	return &repository{
		logger: logger,
		m:      migrator,
		db:     db,
	}
}

func (r *repository) Init(ctx context.Context) error {
	r.logger.DebugContext(ctx, "initializing database")
	return r.m.Init(ctx)
}

func (r *repository) Generate(ctx context.Context, name string) error {

	name = strings.ReplaceAll(name, " ", "_")
	m, err := r.m.CreateGoMigration(ctx, name)
	if err != nil {
		return err
	}

	r.logger.InfoContext(ctx, "created migration file", slog.String("file", m.Name))

	return nil
}

func (r *repository) Migrate(ctx context.Context) error {

	err := r.m.Lock(ctx)
	if err != nil {
		return err
	}
	defer r.m.Unlock(ctx)

	group, err := r.m.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		r.logger.DebugContext(ctx, "there are no new migrations to run (database is up to date)")
		return nil
	}

	r.logger.InfoContext(ctx, "migration complete", slog.String("group", group.Migrations.String()))
	return nil
}

func (r *repository) Rollback(ctx context.Context) error {

	err := r.m.Lock(ctx)
	if err != nil {
		return err
	}
	defer r.m.Unlock(ctx)

	r.logger.InfoContext(ctx, "rolling back migration")
	group, err := r.m.Rollback(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		r.logger.InfoContext(ctx, "there are no migrations to rollback")
		return nil
	}

	r.logger.InfoContext(ctx, "rolled back", slog.String("group", group.String()))
	return nil
}

func (r *repository) Lock(ctx context.Context) error {
	r.logger.InfoContext(ctx, "locking database")
	return r.m.Lock(ctx)
}

func (r *repository) Unlock(ctx context.Context) error {
	r.logger.InfoContext(ctx, "unlocking database")
	return r.m.Unlock(ctx)
}
