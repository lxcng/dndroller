package repo

import (
	"context"
	"database/sql"
	conf "dndroller/internal/config"
	"dndroller/internal/logger"

	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewRepo),
		fx.Invoke(func(lc fx.Lifecycle, r *Repo) {
			lc.Append(fx.Hook{
				OnStart: r.Start,
				OnStop:  r.Stop,
			})
		}),
	)
}

type Repo struct {
	*Client
	cfg *conf.Config
	log *logger.Zap
}

func NewRepo(cfg *conf.Config, log *logger.Zap) *Repo {
	return &Repo{
		cfg: cfg,
		log: log,
	}
}

func (x *Repo) Start(_ context.Context) error {
	client, err := NewCl(x.cfg, x.log)
	if err != nil {
		return err
	}
	x.Client = client
	return nil
}

func (x *Repo) Stop(_ context.Context) error {
	return x.Client.Close()
}

func init() {
	sql.Register("postgres", stdlib.GetDefaultDriver())
}

func NewCl(cfg *conf.Config, log *logger.Zap) (*Client, error) {
	client, err := Open("postgres", cfg.DbUrl, Log(log.Debug))
	if err != nil {
		return nil, log.LogAndWrapError(err, "failed to init ent repo")
	}
	return client, nil
}
