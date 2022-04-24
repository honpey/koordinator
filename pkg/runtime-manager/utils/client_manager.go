package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/koordinator-sh/koordinator/apis/runtime/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"time"

	"github.com/golang/groupcache/lru"
)

type HookServerClientManager struct {
	cache *lru.Cache
}

const (
	defaultCacheSize = 200
)

// TODO: garbage client gc
func NewClientManager() (*HookServerClientManager, error) {
	cache := lru.New(defaultCacheSize)
	return &HookServerClientManager{
		cache: cache,
	}, nil
}

type HookServerPath struct {
	Path string
	Port int64
}

type RuntimeHookClient struct {
	SockPath string
	v1alpha1.RuntimeHookServiceClient
}

func (cm *HookServerClientManager) HookClient(serverPath HookServerPath) (*http.Client, error) {
	cacheKey, err := json.Marshal(serverPath)
	if err != nil {
		return nil, err
	}
	if client, ok := cm.cache.Get(string(cacheKey)); ok {
		return client.(*http.Client), nil
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", serverPath.Path)
			},
		},
		Timeout: 5 * time.Second,
	}
	cm.cache.Add(string(cacheKey), client)
	return client, nil
}

func NewRuntimeHookClient(sockPath string) (*RuntimeHookClient, error) {
	client := &RuntimeHookClient{
		SockPath: sockPath,
	}
	conn, err := grpc.Dial(fmt.Sprintf("unix://%v", sockPath),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		klog.Infof("info %v", err)
		return nil, err
	}
	client.RuntimeHookServiceClient = v1alpha1.NewRuntimeHookServiceClient(conn)
	return client, nil
}

func (cm *HookServerClientManager) RuntimeHookClient(serverPath HookServerPath) (*RuntimeHookClient, error) {
	cacheKey, err := json.Marshal(serverPath)
	if err != nil {
		return nil, err
	}
	if client, ok := cm.cache.Get(string(cacheKey)); ok {
		return client.(*RuntimeHookClient), nil
	}

	runtimeHookClient, err := NewRuntimeHookClient(serverPath.Path)
	if err != nil {
		klog.Infof("err: %v", err)
		return nil, err
	}
	cm.cache.Add(string(cacheKey), runtimeHookClient)
	return runtimeHookClient, nil
}
