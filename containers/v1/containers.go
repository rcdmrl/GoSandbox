package v1

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type DockerCompose struct {
	filePath string
	// env vars:
	pgUser   string
	pgPwd    string
	pgDbName string

	stack compose.ComposeStack
}

func NewDockerCompose(filePath string, pgUser string, pgPwd string, pgDbName string) *DockerCompose {
	return &DockerCompose{filePath, pgUser, pgPwd, pgDbName, nil}
}

func (dc *DockerCompose) getEnvVars() map[string]string {
	m := make(map[string]string)
	m["POSTGRES_USER"] = dc.pgUser
	m["POSTGRES_PASSWORD"] = dc.pgPwd
	m["POSTGRES_DB"] = dc.pgDbName
	return m
}

func (dc *DockerCompose) Start() error {
	stack, err := compose.NewDockerComposeWith(
		compose.StackIdentifier("test"),
		compose.WithStackFiles(dc.filePath),
	)
	if err != nil {
		return err
	}

	dc.stack = compose.ComposeStack(stack)
	dc.stack = dc.stack.WithEnv(dc.getEnvVars())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return dc.stack.Up(ctx)
}

func (dc *DockerCompose) ShutDown() {
	if dc.stack == nil {
		return
	}

	fmt.Println("Shutting down the docker compose stack..")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := dc.stack.Down(
		ctx,
		compose.RemoveOrphans(true),
		compose.RemoveVolumes(false),
		compose.RemoveImagesLocal,
	)
	if err != nil {
		log.Printf("Failed to stop stack: %v", err)
	}
}
