package kube

import (
	"context"
	"fmt"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/tools/portforward"
)

func (c *Client) TunnelingDialerForPod(podName string) (httpstream.Dialer, error) {
	u := c.cs.CoreV1().RESTClient().Get().Resource("pods").
		Namespace(c.Namespace()).Name(podName).SubResource("portforward").URL()
	return c.TunnelingDialer(u)
}

func (c *Client) serviceNameToPod(service string) (string, error) {
	podList, err := c.cs.CoreV1().Pods(c.Namespace()).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + service,
	})
	if err != nil {
		return "", err
	}
	if len(podList.Items) == 0 {
		return "", fmt.Errorf("no pods found for service %s", service)
	}
	return podList.Items[0].Name, nil
}

func (c *Client) PortForward(fromSvc string, listen []string, ports []string) (stopChan, readyChan chan struct{}, outStream, errStream io.ReadCloser, err error) {
	podName, err := c.serviceNameToPod(fromSvc)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	dialer, err := c.TunnelingDialerForPod(podName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	outR, outW := io.Pipe()
	errR, errW := io.Pipe()

	stopChan = make(chan struct{}, 1)
	readyChan = make(chan struct{})
	forwarder, err := portforward.NewOnAddresses(dialer, listen, ports, stopChan, readyChan, outW, errW)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	go forwarder.ForwardPorts()

	return stopChan, readyChan, outR, errR, nil
}
