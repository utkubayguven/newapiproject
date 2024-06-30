package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	client *clientv3.Client
}

func NewEtcdClient(config clientv3.Config) (*EtcdClient, error) {
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &EtcdClient{client: client}, nil
}

func (e *EtcdClient) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) > 0 {
		return resp.Kvs[0].Value, nil
	}
	return nil, nil
}

func (e *EtcdClient) Put(ctx context.Context, key string, value []byte) error {
	_, err := e.client.Put(ctx, key, string(value))
	return err
}

func (e *EtcdClient) Delete(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *EtcdClient) Post(ctx context.Context, key string, value []byte) error {
	_, err := e.client.Put(ctx, key, string(value))
	return err
}

func (e *EtcdClient) Close() error {
	return e.client.Close()
}
