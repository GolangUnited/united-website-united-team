package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	httpV1 "github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1/auth"
	"github.com/zhuravlev-pe/course-watch/internal/app/config"
	"github.com/zhuravlev-pe/course-watch/internal/app/server"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/keygen"
	"github.com/zhuravlev-pe/course-watch/pkg/postgres"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"log"
)

// @title Course Watch API
// @version 1.0
// @description REST API for Course Watch App

// @host localhost:8080
// @BasePath /api/v1/

// @tag.name User
// @tag.description Managing user account

// @tag.name courses
// @tag.description Temporary endpoints for Swagger demo. To be removed

// @tag.name Authentication
// @tag.description Login, logout and other security related operations

// Run initializes whole application.
func Run() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler, cleanup, err := injectHandler(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	srv := server.NewServer(cfg, handler.Init())

	log.Print("Starting server")
	if err = srv.Run(); err != nil {
		log.Fatal(err)
	}
}

// nolint: ireturn
func createIdGen(cfg *config.Config) (service.IdGenerator, error) {
	return idgen.New(cfg.SnowflakeNode)
}

func createPgClient(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, func(), error) {
	pgConfig := postgres.NewPgConfig(
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
	pgClient, err := postgres.NewClient(ctx, pgConfig)
	if err != nil {
		return nil, nil, err
	}
	return pgClient, pgClient.Close, nil
}

// nolint: ireturn
func createAuthenticator(cfg *config.Config) (httpV1.BearerAuthenticator, error) {
	// JwtHandler uses HMAC-SHA256 for signing, block size for SHA256 is 64 bytes, so the key size is the same
	key, err := keygen.Generate(cfg.JWTAuthentication.SigningKey, "bearer-auth.key", 64)
	if err != nil {
		return nil, err
	}
	jwtHandler := security.NewJwtHandler(
		cfg.JWTAuthentication.Issuer,
		cfg.JWTAuthentication.ExpectedAudience,
		cfg.JWTAuthentication.TargetAudience,
		cfg.JWTAuthentication.TokenTTL,
		key,
	)
	bearerAuth := auth.NewBearerAuthenticator(jwtHandler)
	return bearerAuth, nil
}
