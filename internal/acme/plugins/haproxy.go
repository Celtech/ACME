package plugins

import (
	"fmt"
	acmeConfig "github.com/Celtech/ACME/config"
	tcp_socket "github.com/Celtech/ACME/internal/tcp-socket"
	log "github.com/sirupsen/logrus"
	"strings"
)

type HAProxyServer struct {
	Host string
	Port int
}

type HaProxyPlugin struct {
	Enable  bool
	Servers []HAProxyServer
}

func Run(certListPath string, certPath string, contents string) error {
	var haProxyPlugin HaProxyPlugin
	err := acmeConfig.GetConfig().UnmarshalKey("plugins.haproxy", &haProxyPlugin)
	if err != nil {
		return err
	}

	if haProxyPlugin.Enable != true {
		return nil
	}

	for _, server := range haProxyPlugin.Servers {
		log.Infof("Working on server %s:%d", server.Host, server.Port)

		socket := tcp_socket.TCPSocket{
			Address: server.Host,
			Port:    server.Port,
		}

		// Remove empty lines and trim trailing slashes to prevent HAProxy runtime API from yelling
		// at us claiming invalid commands. Blank lines and trailing `\n`'s break things...
		content := strings.TrimSuffix(strings.Replace(contents, "\n\n", "\n", -1), "\n")

		commands := []string{
			"new ssl cert " + certPath,
			fmt.Sprintf("set ssl cert %s <<\n%s\n", certPath, content),
			"commit ssl cert " + certPath,
			fmt.Sprintf("add ssl crt-list %s <<\n%s\n", certListPath, certPath),
		}

		for i, cmd := range commands {
			err := socket.Write(cmd)
			if err != nil {
				log.Errorf("error running command %d: %v", i, err)
				return fmt.Errorf("error running command %d: %v", i, err)
			}
		}
	}

	return nil
}
