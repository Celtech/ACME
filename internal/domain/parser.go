package domain

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func ParseDomainName(domainName string) string {
	url, err := url.Parse(fmt.Sprintf("http://%s", domainName))
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(url.Hostname(), ".")
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
