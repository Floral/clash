package executor

import (
	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	P "github.com/Dreamacro/clash/proxy"
	T "github.com/Dreamacro/clash/tunnel"
)

// Parse config with default config path
func Parse() (*config.Config, error) {
	return ParseWithPath(C.Path.Config())
}

// ParseWithPath parse config with custom config path
func ParseWithPath(path string) (*config.Config, error) {
	return config.Parse(path)
}

// ApplyConfig dispatch configure to all parts
func ApplyConfig(cfg *config.Config, force bool) {
	if force {
		updateGeneral(cfg.General)
	}
	updateProxies(cfg.Proxies)
	updateRules(cfg.Rules)
}

func GetGeneral() *config.General {
	ports := P.GetPorts()
	return &config.General{
		Port:      ports.Port,
		SocksPort: ports.SocksPort,
		RedirPort: ports.RedirPort,
		AllowLan:  P.AllowLan(),
		Mode:      T.Instance().Mode(),
		LogLevel:  log.Level(),
	}
}

func updateProxies(proxies map[string]C.Proxy) {
	T.Instance().UpdateProxies(proxies)
}

func updateRules(rules []C.Rule) {
	T.Instance().UpdateRules(rules)
}

func updateGeneral(general *config.General) {
	allowLan := general.AllowLan

	P.SetAllowLan(allowLan)
	if err := P.ReCreateHTTP(general.Port); err != nil {
		log.Errorln("Start HTTP server error: %s", err.Error())
	}

	if err := P.ReCreateSocks(general.SocksPort); err != nil {
		log.Errorln("Start SOCKS5 server error: %s", err.Error())
	}

	if err := P.ReCreateRedir(general.RedirPort); err != nil {
		log.Errorln("Start Redir server error: %s", err.Error())
	}

	log.SetLevel(general.LogLevel)
	T.Instance().SetMode(general.Mode)
}
