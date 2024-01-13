package config

type NetworkConfig struct {
	IP      string `json:"ip"`
	CIDR    string `json:"cidr"`
	MaxPort int    `json:"max_port"`
	MinPort int    `json:"min_port"`
}

func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		"127.0.0.1", "127.0.0.1/32", 4569, 4565,
	}
}

type ConstNetworkConfig struct {
	c *NetworkConfig
}

func (c *ConstNetworkConfig) GetIP() string   { return c.c.IP }
func (c *ConstNetworkConfig) GetCIDR() string { return c.c.CIDR }
func (c *ConstNetworkConfig) GetMaxPort() int { return c.c.MaxPort }
func (c *ConstNetworkConfig) GetMinPort() int { return c.c.MinPort }

var config = &ConstNetworkConfig{DefaultNetworkConfig()}

func GetConfig() *ConstNetworkConfig { return config }
func SetConfig(nc *NetworkConfig)    { config = &ConstNetworkConfig{nc} }
