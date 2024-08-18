package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	alertmanagerURL   = "http://<PROMETHEUS_ALERTMANAGER_IP>:9093/api/v1/alerts"
	slackWebhookURL   = "https://hooks.slack.com/services/your/slack/webhook/url"
	kubectlCommand    = "/usr/local/bin/kubectl" // Path to kubectl binary
	clusterNamespace  = "default"                // Namespace of the cluster
	podRestartCommand = "kubectl delete pod %s -n %s"
	scaleCommand      = "kubectl scale deployment %s --replicas=%d -n %s"
)

type Alert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname string `json:"alertname"`
		Pod       string `json:"pod"`
	} `json:"labels"`
}

func main() {
	handleAlerts()
}

func handleAlerts() {
	resp, err := http.Get(alertmanagerURL)
	if err != nil {
		log.Fatalf("Failed to connect to Alertmanager: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var alerts []Alert
	err = json.Unmarshal(body, &alerts)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	for _, alert := range alerts {
		if alert.Status == "firing" && alert.Labels.Alertname == "HighCPUUsage" {
			podName := alert.Labels.Pod

			// Restart pod
			err = restartPod(podName, clusterNamespace)
			if err != nil {
				log.Printf("Failed to restart pod: %v", err)
			} else {
				sendSlackNotification(fmt.Sprintf("Pod %s restarted due to high CPU usage", podName))
			}

			// Scale deployment
			err = scaleDeployment("sample-app", 3, clusterNamespace)
			if err != nil {
				log.Printf("Failed to scale deployment: %v", err)
			} else {
				sendSlackNotification(fmt.Sprintf("Scaled deployment sample-app to 3 replicas"))
			}
		}
	}
}

func restartPod(podName string, namespace string) error {
	cmd := fmt.Sprintf(podRestartCommand, podName, namespace)
	return executeCommand(cmd)
}

func scaleDeployment(deploymentName string, replicas int, namespace string) error {
	cmd := fmt.Sprintf(scaleCommand, deploymentName, replicas, namespace)
	return executeCommand(cmd)
}

func executeCommand(command string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed with error: %v, stderr: %v", err, stderr.String())
	}
	return nil
}

func sendSlackNotification(message string) {
	payload := map[string]string{"text": message}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Failed to send Slack notification: %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "ok") {
		log.Printf("Slack notification failed: %s", body)
	}
}
