package mock

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestPostgresqlDB() error {

	cfg := config.ConfigSingleton.GetInstance()
	ctx := context.Background()

	testcontainers.Logger = log.New(io.Discard, "", 0)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": cfg.Database.Password,
			"POSTGRES_USER":     cfg.Database.User,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return err
	}

	cfg.Database.Host, _ = postgresC.Host(ctx)
	natPort, _ := postgresC.MappedPort(ctx, "5432")
	cfg.Database.Port = natPort.Port()

	return nil
}
