package client

import (
	"errors"
	"strings"

	"github.com/lcpu-club/kube-dify-gui/internal/kube"
)

type Client struct {
	addr      string
	token     string
	namespace string

	k *kube.Client
}

func parseNSFromToken(token string) (string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}

	uid := parts[1]
	ns := "u-" + uid

	return ns, nil
}

func New(addr string, token string, ns string) (*Client, error) {
	k, err := kube.NewClient(addr, token, ns)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr:      addr,
		token:     token,
		namespace: ns,
		k:         k,
	}, nil
}

func NewHPCGame(token string) (*Client, error) {
	addr := "https://hpcgame.pku.edu.cn/kube"
	ns, err := parseNSFromToken(token)
	if err != nil {
		return nil, err
	}

	return New(addr, token, ns)
}
