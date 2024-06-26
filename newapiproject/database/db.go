package database

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	Client *clientv3.Client
}

var etcdClient *EtcdClient

func InitEtcd(endpoints []string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	etcdClient = &EtcdClient{Client: client}
	return etcdClient, nil
}

func GetClient(endpoints []string) (*clientv3.Client, error) {
	if etcdClient == nil || etcdClient.Client == nil {
		client, err := InitEtcd(endpoints)
		if err != nil {
			return nil, err
		}
		return client.Client, nil
	}
	return etcdClient.Client, nil
}

// test function
func TestPutGet(client *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Put(ctx, "testkey", "testvalue")
	if err != nil {
		return fmt.Errorf("failed to put value to etcd: %w", err)
	}

	resp, err := client.Get(ctx, "testkey")
	if err != nil {
		return fmt.Errorf("failed to get value from etcd: %w", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	return nil
}
