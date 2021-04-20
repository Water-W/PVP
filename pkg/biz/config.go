package biz

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	MyAddress string
	IsMaster  bool
	Master    struct {
		Address  string
		HTTPPort int
	}
	Worker struct {
		Addresses []string
		URL       string
		NodeQuery string
		LinkQuery string
	}
}

func DefaultConfig() *Config {
	return &Config{
		MyAddress: "127.0.0.1:8001",
		IsMaster:  true,
		Master: struct {
			Address  string
			HTTPPort int
		}{
			Address:  "",
			HTTPPort: -1,
		},
		Worker: struct {
			Addresses []string
			URL       string
			NodeQuery string
			LinkQuery string
		}{
			Addresses: []string{},
			URL:       "http://127.0.0.1:2404/report",
			NodeQuery: "{ID}",
			LinkQuery: "Peers",
		},
	}
}

/*===========================================================================*/

func LoadConfigFromFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("loadConfigFromFile err:%s", err)
	}
	defer f.Close()
	return LoadConfig(f)
}

func LoadConfig(r io.Reader) (*Config, error) {
	c := DefaultConfig()
	dec := toml.NewDecoder(r)
	err := dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) DumpTo(w io.Writer) error {
	enc := toml.NewEncoder(w)
	return enc.Encode(c)
}

func (c *Config) DumpToFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("DumpToFile err:%s", err)
	}
	defer f.Close()
	return c.DumpTo(f)
}

func (c *Config) ToMasterConfig() (*MasterConfig, error) {
	if !c.IsMaster {
		return nil, fmt.Errorf("not a valid master config")
	}
	tcpaddr, err := net.ResolveTCPAddr("tcp", c.MyAddress)
	if err != nil {
		return nil, err
	}
	return &MasterConfig{
		ListenPort: tcpaddr.Port,
	}, nil

}

func (c *Config) ToWorkerConfig() (*WorkerConfig, error) {
	if c.IsMaster {
		return nil, fmt.Errorf("not a valid worker config")
	}
	return &WorkerConfig{
		MasterAddr: c.Master.Address,
		URL:        c.Worker.URL,
		NodeQuery:  c.Worker.NodeQuery,
		LinksQuery: c.Worker.LinkQuery,
	}, nil
}
