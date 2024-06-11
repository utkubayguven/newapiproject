package handlers

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	client *clientv3.Client
}

func NewHandler(client *clientv3.Client) *Handler {
	h := Handler{client: client}
	return &h
}
