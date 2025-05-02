package kube

import (
	"net/url"

	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/discovery"
	cached "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/portforward"
)

const DefaultQPS = 10
const DefaultBurst = 20

type Client struct {
	host      string
	token     string
	namespace string

	rConf *rest.Config

	cs *kubernetes.Clientset
	dc *dynamic.DynamicClient

	mapper *restmapper.DeferredDiscoveryRESTMapper
	cDis   discovery.CachedDiscoveryInterface
}

func NewClient(host string, token string, namespace string) (*Client, error) {
	c := &Client{
		host:      host,
		token:     token,
		namespace: namespace,
	}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) init() error {
	rConf := &rest.Config{
		Host:        c.host,
		BearerToken: c.token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // TODO: set to false
		},
		QPS:   DefaultQPS,
		Burst: DefaultBurst,
	}
	cs, err := kubernetes.NewForConfig(rConf)
	if err != nil {
		return err
	}
	dc, err := dynamic.NewForConfig(rConf)
	if err != nil {
		return err
	}

	cachedClient := cached.NewMemCacheClient(cs.Discovery())
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	c.cs = cs
	c.dc = dc
	c.cDis = cachedClient
	c.mapper = mapper

	c.rConf = rConf

	return nil
}

func (c *Client) Client() *kubernetes.Clientset {
	return c.cs
}

func (c *Client) Dynamic() *dynamic.DynamicClient {
	return c.dc
}

func (c *Client) Host() string {
	return c.host
}

func (c *Client) Namespace() string {
	return c.namespace
}

func (c *Client) TunnelingDialer(u *url.URL) (httpstream.Dialer, error) {
	tunnelingDialer, err := portforward.NewSPDYOverWebsocketDialer(u, c.rConf)
	return tunnelingDialer, err
}
