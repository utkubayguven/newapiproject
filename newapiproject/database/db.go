package database

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	Client *clientv3.Client
}

// InitEtcd initializes the etcd connection
func InitEtcd(endpoints []string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://etcd1:2379", "http://etcd2:2378", "http://etcd3:2377"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}

	return &EtcdClient{Client: client}, nil
}
