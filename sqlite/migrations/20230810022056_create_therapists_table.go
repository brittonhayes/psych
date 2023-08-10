package migrations

import (
	"context"

	"github.com/brittonhayes/therapy"
	"github.com/uptrace/bun"
)

var models = []interface{}{
	(*therapy.Therapist)(nil),
}

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		for _, model := range models {
			_, err := db.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		for _, model := range models {
			_, err := db.NewDropTable().IfExists().Model(model).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
