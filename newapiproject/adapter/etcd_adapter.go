package adapter

import (
	"context"

	"newapiprojet/database"
	"newapiprojet/etcd"
)

type EtcdAdapter struct {
	client *etcd.EtcdClient
}

func NewEtcdAdapter(client *etcd.EtcdClient) database.Database {
	return &EtcdAdapter{client: client}
}

func (e *EtcdAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	return e.client.Get(ctx, key)
}

func (e *EtcdAdapter) Put(ctx context.Context, key string, value []byte) error {
	return e.client.Put(ctx, key, value)
}

func (e *EtcdAdapter) Delete(ctx context.Context, key string) error {
	return e.client.Delete(ctx, key)
}

func (e *EtcdAdapter) Post(ctx context.Context, key string, value []byte) error {
	return e.client.Post(ctx, key, value)
}
