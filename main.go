package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
)

var (
	config             SlackproxyConfig
	SETTING_JSON       = os.Getenv("SETTING_JSON")
	defaultSettingJson = "setting.json"
)

type OutgoingContent struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func processProxy(token string, user_name string, text string) {

	cfg, ok := config.ProxyList[token]
	if !ok {
		glog.Error("does not exists %s", token)
		return
	}

	msg := fmt.Sprintf("(%s) %s", cfg.FromNetworkName, text)
	data := OutgoingContent{user_name, msg}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		glog.Error("json.Marshal failed")
		return
	}

	req, err := http.NewRequest(
		"POST",
		cfg.OutgoingUrl,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		glog.Error("http.NewRequset failed")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(15 * time.Second)}
	resp, err := client.Do(req)
	defer resp.Body.Close()
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := r.FormValue("token")
	username := r.FormValue("user_name")
	text := r.FormValue("text")

	processProxy(token, username, text)

	fmt.Fprintf(w, "{}")
}

func main() {
	if SETTING_JSON == "" {
		SETTING_JSON = defaultSettingJson
	}
	var err error
	config, err = LoadSetting(SETTING_JSON)
	if err != nil {
		glog.Errorf("%#v", err)
		panic("failed to load config")
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/post", handlePost)
	http.ListenAndServe(config.Bind, nil)
}
