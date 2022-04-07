package domain

import (
	"fmt"
	"regexp"
)

func ParseDomainName(domainName string) (string, error) {
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/?\n]+)`)

	found := re.FindAllString(domainName, -1)
	if len(found) > 0 {
		return found[0], nil
	}

	return "", fmt.Errorf("failed to parse host from %s, a valid domain name must be supplied", domainName)
}
