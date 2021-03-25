package etcd

import (
	"context"
	"fmt"
	"github.com/tangxusc/crypt/backend"

	goetcd "go.etcd.io/etcd/client/v3"
)

type Client struct {
	client    *goetcd.Client
	keysAPI   goetcd.KV
	waitIndex uint64
}

func New(machines []string) (*Client, error) {
	newClient, err := goetcd.New(goetcd.Config{
		Endpoints: machines,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new etcd client for crypt.backend.Client: %v", err)
	}
	keysAPI := goetcd.NewKV(newClient)
	return &Client{client: newClient, keysAPI: keysAPI, waitIndex: 0}, nil
}

func (c *Client) Get(key string) ([]byte, error) {
	return c.GetWithContext(context.TODO(), key)
}

func (c *Client) GetWithContext(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.keysAPI.Get(ctx, key, nil)
	if err != nil {
		return nil, err
	}
	return resp.Kvs[0].Value, nil
}

func (c *Client) List(key string) (backend.KVPairs, error) {
	return c.ListWithContext(context.TODO(), key)
}

func (c *Client) ListWithContext(ctx context.Context, key string) (backend.KVPairs, error) {
	resp, err := c.keysAPI.Get(ctx, key, nil)
	if err != nil {
		return nil, err
	}
	pairs := make([]*backend.KVPair, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		pairs[i] = &backend.KVPair{
			Key:   string(kv.Key),
			Value: kv.Value,
		}
	}
	return pairs, nil
}

func (c *Client) Set(key string, value []byte) error {
	return c.SetWithContext(context.TODO(), key, value)
}

func (c *Client) SetWithContext(ctx context.Context, key string, value []byte) error {
	_, err := c.keysAPI.Put(ctx, key, string(value), nil)
	return err
}

func (c *Client) Watch(key string, stop chan bool) <-chan *backend.Response {
	return c.WatchWithContext(context.Background(), key, stop)
}

func (c *Client) WatchWithContext(ctx context.Context, key string, stop chan bool) <-chan *backend.Response {
	subCtx, cancel := context.WithCancel(ctx)
	respChan := make(chan *backend.Response, 0)
	go func() {
		watch := c.client.Watch(subCtx, key)
		for {
			select {
			case <-subCtx.Done():
				cancel()
			case <-stop:
				cancel()
			case <-watch:
				get, err := c.keysAPI.Get(subCtx, key)
				if err != nil {
					respChan <- &backend.Response{Error: err}
				}
				respChan <- &backend.Response{Value: get.Kvs[0].Value}
			}
		}
	}()
	return respChan
}
