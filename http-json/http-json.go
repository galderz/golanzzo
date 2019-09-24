package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	host := "192.168.99.137:30831"
	fmt.Printf("Exposed service host is %v\n", host)
	client := &http.Client{}
	//putViaRoute(host, "test-operator", client)
	resp := getViaRoute(host, client)
	fmt.Printf("Response: %s\n", resp)
	view := clusterView(host, client)
	fmt.Printf("Cluster view: %s\n", view)
}

type ClusterHealth struct {
	Nodes []string `json:"node_names"`
}

type Health struct {
	ClusterHealth  ClusterHealth `json:"cluster_health"`
}

func clusterView(host string, client *http.Client) []string {
	req, _ := http.NewRequest("GET", "http://"+host+"/rest/v2/cache-managers/DefaultCacheManager/health", nil)
	req.SetBasicAuth("developer", "mjLMCQdMOUBtkmyC")
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected response %v", resp))
	}
	var health Health
	err = json.NewDecoder(resp.Body).Decode(&health)
	if err != nil {
		panic(fmt.Errorf("unable to decode"))
	}
	return health.ClusterHealth.Nodes
}

func getViaRoute(host string, client *http.Client) string {
	req, _ := http.NewRequest("GET", "http://"+host+"/rest/v2/cache-managers/DefaultCacheManager/health", nil)
	req.SetBasicAuth("developer", "mjLMCQdMOUBtkmyC")
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("unexpected response %v", resp))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(bodyBytes)
}
