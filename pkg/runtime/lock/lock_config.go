package lock

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	strategyKey = "keyPrefix"

	strategyAppid     = "appid"
	strategyStoreName = "name"
	strategyNone      = "none"
	strategyDefault   = strategyAppid

	apiPrefix    = "lock"
	apiSeparator = "|||"
	separator    = "||"
)

var lockConfiguration = map[string]*StoreConfiguration{}

type StoreConfiguration struct {
	keyPrefixStrategy string
}

func SaveLockConfiguration(storeName string, metadata map[string]string) error {
	strategy := strings.ToLower(metadata[strategyKey])
	if strategy == "" {
		strategy = strategyDefault
	} else {
		err := checkKeyIllegal(metadata[strategyKey])
		if err != nil {
			return err
		}
	}

	lockConfiguration[storeName] = &StoreConfiguration{keyPrefixStrategy: strategy}
	return nil
}

func GetModifiedLockKey(key, storeName, appID string) (string, error) {
	if err := checkKeyIllegal(key); err != nil {
		return "", err
	}
	config := getConfiguration(storeName)
	switch config.keyPrefixStrategy {
	case strategyNone:
		return fmt.Sprintf("%s%s%s", apiPrefix, apiSeparator, key), nil
	case strategyStoreName:
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, storeName, separator, key), nil
	case strategyAppid:
		if appID == "" {
			return fmt.Sprintf("%s%s%s", apiPrefix, apiSeparator, key), nil
		}
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, appID, separator, key), nil
	default:
		return fmt.Sprintf("%s%s%s%s%s", apiPrefix, apiSeparator, config.keyPrefixStrategy, separator, key), nil
	}
}

func getConfiguration(storeName string) *StoreConfiguration {
	c := lockConfiguration[storeName]
	if c == nil {
		c = &StoreConfiguration{keyPrefixStrategy: strategyDefault}
		lockConfiguration[storeName] = c
	}

	return c
}

func checkKeyIllegal(key string) error {
	if strings.Contains(key, separator) {
		return errors.Errorf("input key/keyPrefix '%s' can't contain '%s'", key, separator)
	}
	return nil
}
