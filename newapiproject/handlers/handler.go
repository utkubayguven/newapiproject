package handlers

import (
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	client *clientv3.Client
}

func NewHandler(client *clientv3.Client) *Handler {
	if client == nil {
		log.Fatalf("etcd client is nil")
	}
	return &Handler{client: client}
}
