package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	host := "cache-infinispan-0-namespace-for-testing.router.172.17.0.3.nip.io"
	fmt.Printf("Exposed service host is %v\n", host)
	client := &http.Client{}
	//putViaRoute(host, "test-operator", client)
	getViaRoute(host, client)
}

func getViaRoute(host string, client *http.Client) string {
	req, _ := http.NewRequest("GET", "http://"+host+"/rest/default/test", nil)
	req.SetBasicAuth("infinispan", "infinispan")
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

//func putViaRoute(host string, value string, client *http.Client) {
//	body := bytes.NewBuffer([]byte(value))
//	req, err := http.NewRequest("POST", "http://"+host+"/rest/default/test", body)
//	req.Header.Set("Content-Type", "text/plain")
//	req.SetBasicAuth("infinispan", "infinispan")
//	fmt.Printf("Put request via route: %v\n", req)
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err.Error())
//	}
//	if resp.StatusCode != http.StatusOK {
//		panic(fmt.Errorf("unexpected response %v", resp))
//	}
//}
