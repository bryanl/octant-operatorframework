package oof

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/bryanl/octant-operatorframework/thirdparty/operator-registry/pkg/api"
)

//go:generate mockgen -destination=./fake/mock_registry_client.go -package=fake github.com/bryanl/octant-operatorframework/pkg/oof RegistryClient

type RegistryClient interface {
	ListPackages(ctx context.Context) ([]string, error)
	Close() error
}

type GRPCRegistryClient struct {
	conn *grpc.ClientConn
}

var _ RegistryClient = (*GRPCRegistryClient)(nil)

func NewGRPCRegistryClient(address string) (RegistryClient, error) {
	logger.Printf("connecting to registry at %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	client := GRPCRegistryClient{
		conn: conn,
	}

	return &client, nil
}

func (c *GRPCRegistryClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCRegistryClient) ListPackages(ctx context.Context) ([]string, error) {
	client := api.NewRegistryClient(c.conn)

	pkgStream, err := client.ListPackages(ctx, &api.ListPackageRequest{})
	if err != nil {
		return nil, err
	}

	var list []string

	for {
		pkgName, err := pkgStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("list packages: %w", err)
		}

		list = append(list, pkgName.Name)
	}

	return list, nil
}
