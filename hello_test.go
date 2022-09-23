package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	hellov1 "github.com/echovl/hello-grpc/gen/proto/hello/v1"
	"google.golang.org/grpc"
)

const testServerAddr = ":3000"

type helloRespCtxKey struct{}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func callHelloWithUsername(ctx context.Context, username string) (context.Context, error) {
	client, err := newHelloServiceClient()
	if err != nil {
		return ctx, fmt.Errorf("creating client: %s", err)
	}

	resp, err := client.Hello(context.Background(), &hellov1.HelloRequest{Username: username})
	if err != nil {
		return ctx, fmt.Errorf("got service error: %s", err)
	}
	return context.WithValue(ctx, helloRespCtxKey{}, resp.Msg), nil
}

func helloServiceIsRunning() error {
	go runHelloServiceServer()
	return nil
}

func helloShouldReturnMessage(ctx context.Context, wantMsg string) error {
	gotMsg, ok := ctx.Value(helloRespCtxKey{}).(string)
	if !ok {
		return errors.New("no hello response")
	}

	if gotMsg != wantMsg {
		return fmt.Errorf("expected %s response, but got %s", wantMsg, gotMsg)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^call Hello with username "([^"]*)"$`, callHelloWithUsername)
	ctx.Step(`^Hello service is running$`, helloServiceIsRunning)
	ctx.Step(`^Hello should return message "([^"]*)"$`, helloShouldReturnMessage)
}

func runHelloServiceServer() {
	NewServer(testServerAddr).ListenAndServe()
}

func newHelloServiceClient() (hellov1.HelloServiceClient, error) {
	conn, err := grpc.Dial(testServerAddr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return hellov1.NewHelloServiceClient(conn), nil
}
