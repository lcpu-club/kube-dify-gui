package client

import (
	"context"
	"fmt"
	"io"
)

func (c *Client) DoInit() (initPassword string, err error) {
	secYml, iPwd := c.secretsYaml()
	pvcYml := c.pvcYaml()

	err = c.k.Create(context.Background(), secYml, false)
	if err != nil {
		return "", err
	}

	err = c.k.Create(context.Background(), pvcYml, false)
	if err != nil {
		return "", err
	}

	return iPwd, nil
}

func (c *Client) DoStart() error {
	yml := c.deployYaml()
	err := c.k.Create(context.Background(), yml, false)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DoStop() error {
	yml := c.deployYaml()
	err := c.k.Delete(context.Background(), yml, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DoDelete() error {
	c.DoStop()

	yml := c.deleteYaml()
	err := c.k.Delete(context.Background(), yml, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DoPortForward(addr string, port string) (stopChan chan struct{}, outStream, errStream io.ReadCloser, err error) {
	realPortStr := fmt.Sprintf("%s:80", port)

	stopChan, readyChan, outStream, errStream, err := c.k.PortForward("dify-nginx", []string{addr}, []string{realPortStr})
	if err != nil {
		return nil, nil, nil, err
	}
	_ = readyChan

	return stopChan, outStream, errStream, nil
}
