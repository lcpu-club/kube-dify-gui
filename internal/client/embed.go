package client

import (
	_ "embed"
	"encoding/base64"
	"math/rand"
	"strings"
)

//go:embed deploy.yaml
var deployYaml string

//go:embed pvc.yaml
var pvcYaml string

//go:embed secrets.yaml
var secretsYaml string

//go:embed delete.yaml
var deleteYaml string

func replaceNamespace(yaml string, ns string) string {
	return strings.ReplaceAll(yaml, "${.Namespace}", ns)
}

func replaceSecrets(yaml string, secrets map[string]string) string {
	for k, v := range secrets {
		yaml = strings.ReplaceAll(yaml, "${."+k+"}", v)
	}
	return yaml
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateSecrets() (secs map[string]string, initPwd string) {
	secs = make(map[string]string)

	initPwd = randString(10)
	pgPwd := randString(16)
	redisPwd := randString(16)

	secs["PGPassword"] = base64.StdEncoding.EncodeToString([]byte(pgPwd))
	secs["RedisPassword"] = base64.StdEncoding.EncodeToString([]byte(redisPwd))
	secs["InitPassword"] = base64.StdEncoding.EncodeToString([]byte(initPwd))
	return
}

func (c *Client) deployYaml() string {
	return replaceNamespace(deployYaml, c.namespace)
}

func (c *Client) pvcYaml() string {
	return replaceNamespace(pvcYaml, c.namespace)
}

func (c *Client) secretsYaml() (yml string, initPwd string) {
	secs, initPwd := generateSecrets()
	yml = replaceNamespace(secretsYaml, c.namespace)
	yml = replaceSecrets(yml, secs)
	return yml, initPwd
}

func (c *Client) deleteYaml() string {
	return replaceNamespace(deleteYaml, c.namespace)
}
