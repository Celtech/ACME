package configFactory

import (
	"strconv"
)

type serverConfig struct {
	Host string
	Port int
}

type dnsConfig struct {
	Provider string
	Timeout  int
}

type acmeProvider struct {
	ProductionProvider string
	StagingProvider    string
	Email              string
}

type ConfigFactory struct {
	DebugMode           bool
	DataPath            string
	Server              *serverConfig
	HttpChallengeServer *serverConfig
	TLSChallengeServer  *serverConfig
	DNSChallengeServer  *dnsConfig
	AcmeProvider        *acmeProvider
}

func NewConfigFactory() *ConfigFactory {
	configDebugMode, _ := strconv.ParseBool(
		VariableToSetting(Variable{
			name:         "BAKER_DEBUG",
			defaultValue: "false",
			required:     false,
		}))

	configAcmeProviderEmail := VariableToSetting(Variable{
		name:         "BAKER_ACME_PROVIDER_EMAIL",
		defaultValue: "",
		required:     true,
	})

	configAcmeProviderProduction := VariableToSetting(Variable{
		name:         "BAKER_ACME_PROVIDER_HOST",
		defaultValue: "https://acme-v02.api.letsencrypt.org/directory",
		required:     false,
	})

	configAcmeProviderStaging := VariableToSetting(Variable{
		name:         "BAKER_ACME_PROVIDER_STAGING_HOST",
		defaultValue: "https://acme-staging-v02.api.letsencrypt.org/directory",
		required:     false,
	})

	configDataPath := VariableToSetting(Variable{
		name:         "BAKER_DATA_PATH",
		defaultValue: "/data",
		required:     false,
	})

	configServerHost := VariableToSetting(Variable{
		name:         "BAKER_SERVER_HOST",
		defaultValue: "0.0.0.0",
		required:     false,
	})

	configServerPort, _ := strconv.Atoi(VariableToSetting(Variable{
		name:         "BAKER_SERVER_PORT",
		defaultValue: "9022",
		required:     false,
	}))

	configHTTPChallengeServerHost := VariableToSetting(Variable{
		name:         "BAKER_HTTP_CHALLENGE_SERVER_HOST",
		defaultValue: "0.0.0.0",
		required:     false,
	})

	configHTTPChallengeServerPort, _ := strconv.Atoi(VariableToSetting(Variable{
		name:         "BAKER_HTTP_CHALLENGE_SERVER_PORT",
		defaultValue: "80",
		required:     false,
	}))

	configTLSChallengeServerHost := VariableToSetting(Variable{
		name:         "BAKER_TLS_CHALLENGE_SERVER_HOST",
		defaultValue: "0.0.0.0",
		required:     false,
	})

	configTLSChallengeServerPort, _ := strconv.Atoi(VariableToSetting(Variable{
		name:         "BAKER_TLS_CHALLENGE_SERVER_PORT",
		defaultValue: "443",
		required:     false,
	}))

	configDNSChallengeProvider := VariableToSetting(Variable{
		name:         "BAKER_DNS_CHALLENGE_PROVIDER",
		defaultValue: "dnsmadeeasy",
		required:     false,
	})

	return &ConfigFactory{
		DebugMode: configDebugMode,
		DataPath:  configDataPath,
		AcmeProvider: &acmeProvider{
			Email:              configAcmeProviderEmail,
			ProductionProvider: configAcmeProviderProduction,
			StagingProvider:    configAcmeProviderStaging,
		},
		Server: &serverConfig{
			Host: configServerHost,
			Port: configServerPort,
		},
		HttpChallengeServer: &serverConfig{
			Host: configHTTPChallengeServerHost,
			Port: configHTTPChallengeServerPort,
		},
		TLSChallengeServer: &serverConfig{
			Host: configTLSChallengeServerHost,
			Port: configTLSChallengeServerPort,
		},
		DNSChallengeServer: &dnsConfig{
			Provider: configDNSChallengeProvider,
		},
	}
}
