package handlers

import (
	"newapiprojet/database"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Handler struct {
	endpoints []string
}

func NewHandler(endpoints []string) *Handler {
	return &Handler{endpoints: endpoints}
}

func (h *Handler) getClient() (*clientv3.Client, error) {
	return database.GetClient(h.endpoints)
}
