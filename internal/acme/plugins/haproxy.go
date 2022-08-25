package plugins

import (
	"fmt"
	acmeConfig "github.com/Celtech/ACME/config"
	tcp_socket "github.com/Celtech/ACME/internal/tcp-socket"
	"strings"
)

func Run(certListPath string, certPath string, contents string) error {
	enable := acmeConfig.GetConfig().GetBool("plugins.haproxy.enable")
	if enable != true {
		return nil
	}

	host := acmeConfig.GetConfig().GetString("plugins.haproxy.host")
	port := acmeConfig.GetConfig().GetInt("plugins.haproxy.port")

	socket := tcp_socket.TCPSocket{
		Address: host,
		Port:    port,
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
			return fmt.Errorf("error running command %d: %v", i, err)
		}
	}

	return nil
}
