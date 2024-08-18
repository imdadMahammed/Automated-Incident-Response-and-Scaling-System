import os
import requests
from kubernetes import client, config

# Load kube config
config.load_kube_config()

# Set up Kubernetes API client
v1 = client.CoreV1Api()
apps_v1 = client.AppsV1Api()

# URL for Prometheus Alertmanager API
ALERTMANAGER_URL = "http://<PROMETHEUS_ALERTMANAGER_IP>:9093/api/v1/alerts"

# Slack webhook URL
SLACK_WEBHOOK_URL = "https://hooks.slack.com/services/your/webhook/url"

# Function to restart a pod
def restart_pod(namespace, pod_name):
    print(f"Restarting pod: {pod_name}")
    try:
        v1.delete_namespaced_pod(name=pod_name, namespace=namespace)
        print(f"Pod {pod_name} in namespace {namespace} successfully restarted.")
        send_slack_notification(f"Pod {pod_name} in namespace {namespace} successfully restarted.")
    except client.exceptions.ApiException as e:
        print(f"Exception when restarting pod: {e}")
        send_slack_notification(f"Failed to restart pod {pod_name}: {e}")

# Function to scale a deployment
def scale_deployment(namespace, deployment_name, replicas):
    print(f"Scaling deployment: {deployment_name} to {replicas} replicas.")
    try:
        scaling = apps_v1.read_namespaced_deployment(deployment_name, namespace)
        scaling.spec.replicas = replicas
        apps_v1.replace_namespaced_deployment(deployment_name, namespace, scaling)
        print(f"Scaled {deployment_name} to {replicas} replicas.")
        send_slack_notification(f"Scaled {deployment_name} to {replicas} replicas.")
    except client.exceptions.ApiException as e:
        print(f"Exception when scaling deployment: {e}")
        send_slack_notification(f"Failed to scale deployment {deployment_name}: {e}")

# Function to send a notification to Slack
def send_slack_notification(message):
    payload = {
        "text": message
    }
    try:
        requests.post(SLACK_WEBHOOK_URL, json=payload)
    except requests.exceptions.RequestException as e:
        print(f"Failed to send Slack notification: {e}")

# Function to handle alerts from Prometheus
def handle_alerts():
    try:
        response = requests.get(ALERTMANAGER_URL)
        alerts = response.json()['data']
        for alert in alerts:
            if alert['status'] == 'firing' and alert['labels']['alertname'] == 'HighCPUUsage':
                pod_name = alert['labels']['pod']
                namespace = 'default'

                # Restart the pod with high CPU usage
                restart_pod(namespace, pod_name)

                # Auto-scale the deployment (scale up by 1 replica)
                deployment_name = "sample-app"  # Replace with your deployment name
                scale_deployment(namespace, deployment_name, replicas=3)
                
    except requests.exceptions.RequestException as e:
        print(f"Failed to connect to Alertmanager: {e}")
    except KeyError as e:
        print(f"Key error: {e}")

if __name__ == "__main__":
    handle_alerts()
