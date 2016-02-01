package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

type SlackproxyEntry struct {
	OutgoingUrl     string `json:"url"`
	FromNetworkName string `json:"net"`
}

type SlackproxyConfig struct {
	ProxyList map[string]SlackproxyEntry `json:"proxy_list"`
	Bind      string                     `json:"bind"`
}

func LoadSetting(settingFile string) (SlackproxyConfig, error) {
	var d SlackproxyConfig
	jstring, err := ioutil.ReadFile(settingFile)
	if err != nil {
		glog.Error(err)
		return d, err
	}
	err = json.Unmarshal(jstring, &d)
	if err != nil {
		glog.Error(err)
		return d, err
	}
	return d, nil
}
