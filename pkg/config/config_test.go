package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig(".", "test_config")
	assert.Equal(t, ":8088", Conf.Port)
	assert.Equal(t, "fasthttp", Conf.ServerType)
	assert.Equal(t, 2, len(Conf.Peers))
	assert.Equal(t, 1, len(Conf.Orderers))
	assert.Equal(t, "peer0.org1.example.com:7051", Conf.Peers[0].Addr)
	assert.Equal(t, "/testkey.pem", Conf.Peers[0].Key)
	assert.Equal(t, "/testcert.pem", Conf.Peers[0].Cert)
	assert.Equal(t, "peer0.org1.example.com", Conf.Peers[0].OverrideHostname)
	assert.Equal(t, 1, len(Conf.Peers[0].RootCerts))
	assert.Equal(t, "/testroot.pem", Conf.Peers[0].RootCerts[0])
}
