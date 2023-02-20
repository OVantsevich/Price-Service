package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

const (
	testRedisPort = "9999"
	testRedisHost = "localhost"

	testRedisStream = "testStream"
)

var redisRps *Redis

var streamPool *StreamPool

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
			"--save 20 1",
			"--loglevel warning",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379/tcp": {{HostIP: testRedisHost, HostPort: fmt.Sprintf("%s/tcp", testRedisPort)}},
		}})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	ctx := context.Background()
	if err = pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", testRedisHost, testRedisPort),
		})
		_, err = client.Ping(ctx).Result()
		if err != nil {
			return err
		}
		redisRps = NewRedis(client, testRedisStream)
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to redis: %s", err)
	}
	streamPool = NewStreamPool()

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
