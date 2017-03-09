package config

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const Version = "1"

type Target struct {
	Template      string                    `json:"template"`
	Notifications map[string][]Notification `json:"notifications,omitempty"`
}

type Notification map[string]string

func GetTarget(name string) (*Target, error) {
	targets := viper.GetStringMap("targets")

	for k, t := range targets {
		if k == name {
			target, err := castToTargetType(t)
			if err != nil {
				return nil, fmt.Errorf("cannot cast target %s to Target struct", name)
			}
			return target, nil
		}
	}
	return nil, fmt.Errorf("could not find target %s", name)
}

// TODO rz - UGhh gross
func castToTargetType(t interface{}) (*Target, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling target to intermediary json")
	}
	target := Target{}
	if err := json.Unmarshal(b, &target); err != nil {
		return nil, errors.Wrap(err, "unmarshaling target intemediary json")
	}
	return &target, nil
}
