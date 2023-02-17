package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

const (
	testRedisPort = "99999"
)

var redis *Redis

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("unix:///home/olegvantsevich/.docker/desktop/docker.sock")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "latest",
		Cmd: []string{
			"-save 20 1",
			"-loglevel=warning",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: "localhost", HostPort: fmt.Sprintf("%s/tcp", testRedisPort)}},
		}})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	ctx := context.Background()
	if err := pool.Retry(func() error {
		var err error
		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		})
		defer client.Close()
		db, err = pgxpool.New(ctx, fmt.Sprintf("postgres://postgres:postgres@localhost:%s/userService?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping(ctx)
	}); err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
