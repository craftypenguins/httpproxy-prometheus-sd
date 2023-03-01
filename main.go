package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ScrapeTarget struct {
	Targets []string               `json:"targets"`
	Labels  map[string]interface{} `json:"labels,omitempty"`
}

type HttpProxyList struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				BlackboxMonitor string `json:"blackbox-monitor"`
			} `json:"annotations"`
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
		Spec struct {
			Virtualhost struct {
				Fqdn string `json:"fqdn"`
				TLS  struct {
					SecretName string `json:"secretName"`
				} `json:"tls"`
			} `json:"virtualhost"`
		} `json:"spec"`
	} `json:"items"`
	Kind string `json:"kind"`
}

func getToken() string {
	// tokenPath := filepath.Join("/home/develop/crafty/software/httpproxy-exporter", "sample.token")
	tokenPath := filepath.Join("/var/run/secrets/kubernetes.io/serviceaccount", "token")

	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading token file: %v\n", err)
		os.Exit(1)
	}

	token := string(tokenBytes)
	return token
}

func getProxies() HttpProxyList {

	client := &http.Client{}
	token := getToken()

	//req, err := http.NewRequest("GET", "http://localhost:8000/sample.json", nil)
	req, err := http.NewRequest("GET", "https://kubernetes.default.svc/api/v1/namespaces/default/pods", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}

	//fmt.Fprintf(os.Stderr, "Token: %s", token)
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error making request: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response body: %v\n", err)
		os.Exit(1)
	}

	//body := string(bodyBytes)
	var proxyList HttpProxyList
	err = json.Unmarshal(bodyBytes, &proxyList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling response body: %v\n", err)
		os.Exit(1)
	}

	return proxyList
}

func generateTargets() ([]byte, error) {

	proxyList := getProxies()
	scrapeTargets := [1]ScrapeTarget{}

	for _, proxy := range proxyList.Items {
		if proxy.Metadata.Annotations.BlackboxMonitor == "true" {
			fmt.Println(proxy.Spec.Virtualhost.Fqdn)
			var url string
			if proxy.Spec.Virtualhost.TLS.SecretName != "" {
				url = "https://" + proxy.Spec.Virtualhost.Fqdn
			} else {
				url = "http://" + proxy.Spec.Virtualhost.Fqdn
			}
			scrapeTargets[0].Targets = append(scrapeTargets[0].Targets, url)
		}
	}

	jsonData, err := json.Marshal(scrapeTargets)
	return jsonData, err
}

func prometheusTargets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := generateTargets()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating targets: %v\n", err)
		os.Exit(1)
	}
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/", prometheusTargets)
	log.Fatal(http.ListenAndServe(":8001", nil))
}
