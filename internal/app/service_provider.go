package app

import (
	"context"
	"log"

	"github.com/IRXCI/auth/internal/api/handler"
	"github.com/IRXCI/auth/internal/client/db"
	"github.com/IRXCI/auth/internal/client/db/pg"
	"github.com/IRXCI/auth/internal/client/db/transaction"
	"github.com/IRXCI/auth/internal/closer"
	"github.com/IRXCI/auth/internal/config"
	"github.com/IRXCI/auth/internal/repository"
	authRepository "github.com/IRXCI/auth/internal/repository/auth"
	"github.com/IRXCI/auth/internal/service"
	authService "github.com/IRXCI/auth/internal/service/auth"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.AuthRepository

	auth service.AuthService

	handler *handler.Implementation
}

func newServProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.AuthRepository {
	if s.noteRepository == nil {
		s.noteRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) service.AuthService {
	if s.auth == nil {
		s.auth = authService.NewService(
			s.NoteRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.auth
}

func (s *serviceProvider) NoteImpl(ctx context.Context) *handler.Implementation {
	if s.handler == nil {
		s.handler = handler.NewImplementation(s.NoteService(ctx))
	}

	return s.handler
}
