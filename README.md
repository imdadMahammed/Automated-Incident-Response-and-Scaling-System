
# Kubernetes Incident Responder and Auto-Scaler

This Golang program automates incident response and scaling in a Kubernetes cluster. It integrates with Prometheus Alertmanager to monitor alerts, restarts pods when necessary, scales deployments based on conditions, and sends notifications to a Slack channel for every incident handled.

## Features

- **Alert Monitoring**: Listens to Prometheus Alertmanager for alerts related to high CPU usage.
- **Automated Incident Response**: Automatically restarts failing pods in response to high CPU usage alerts.
- **Auto-Scaling**: Scales Kubernetes deployments based on predefined conditions.
- **Slack Notifications**: Sends notifications to a Slack channel for every pod restart and scaling action.

## Prerequisites

- **Golang**: You must have Golang installed to build and run this program. You can install it [here](https://golang.org/dl/).
- **Kubernetes Cluster**: Ensure you have a running Kubernetes cluster with the appropriate configurations.
- **Kubectl**: The Kubernetes CLI tool `kubectl` must be installed and configured to communicate with your cluster. You can install it [here](https://kubernetes.io/docs/tasks/tools/).
- **Prometheus & Alertmanager**: Prometheus must be installed and monitoring your Kubernetes cluster, and Alertmanager must be configured to trigger alerts.
- **Slack Webhook**: You need a Slack Webhook URL to send notifications to a specific Slack channel. You can create one [here](https://api.slack.com/messaging/webhooks).

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/k8s-incident-responder.git
   cd k8s-incident-responder
   ```

2. **Update Configurations**:
   - Open the `main.go` file and replace the following placeholders with your actual values:
     - `<PROMETHEUS_ALERTMANAGER_IP>`: The IP or URL of your Prometheus Alertmanager instance.
     - `https://hooks.slack.com/services/your/slack/webhook/url`: Your Slack webhook URL.

3. **Build the Application**:
   ```bash
   go build -o k8s-incident-responder main.go
   ```

4. **Run the Application**:
   ```bash
   ./k8s-incident-responder
   ```

## Usage

The program automatically fetches alerts from Prometheus Alertmanager and triggers incident responses based on the alert type. Specifically, it will:

- **Restart Pods**: In the case of high CPU usage, it will automatically restart the failing pod using `kubectl delete pod`.
- **Scale Deployments**: It will scale up the deployment to 3 replicas (can be adjusted in the code) when high CPU usage is detected.
- **Send Slack Notifications**: After handling incidents, it will send notifications to a Slack channel via the webhook.

### Example Output

When a pod is restarted:
```
Restarting pod: my-pod-name
Pod my-pod-name in namespace default successfully restarted.
```

When a deployment is scaled:
```
Scaling deployment: sample-app to 3 replicas.
Scaled sample-app to 3 replicas.
```

### Automating with Cron Jobs

To continuously monitor alerts and respond to incidents, you can automate the program by setting up a cron job:

1. Open the cron editor:
   ```bash
   crontab -e
   ```

2. Add a cron job to run the program every 5 minutes:
   ```bash
   */5 * * * * /path/to/k8s-incident-responder >> /var/log/k8s-incident-responder.log 2>&1
   ```

## Customization

### Changing Deployment Scaling Logic

To customize the scaling logic (e.g., scale up/down based on conditions):
- Modify the `scaleDeployment` function inside the `main.go` file.
- Adjust the number of replicas or implement more complex scaling algorithms.

### Handling Different Alerts

To handle different types of alerts (e.g., memory usage, disk space issues):
- Update the `handleAlerts()` function to listen for additional alert types by modifying the conditions in the alert check logic.

## Configuration Details

- **Alertmanager URL**: The program queries the Alertmanager API for firing alerts.
- **Slack Webhook**: Integrates with Slack to send incident notifications.
- **Kubectl**: Uses `kubectl` commands to interact with the Kubernetes cluster for pod restarts and deployment scaling.

## Error Handling

The program includes basic error handling for:
- HTTP request failures when connecting to Prometheus Alertmanager.
- Command execution failures when restarting pods or scaling deployments.
- Slack notification delivery errors.

Any errors encountered during execution are logged in the console.

## Example Alertmanager Rule

Here’s an example Prometheus alerting rule for high CPU usage:

```yaml
groups:
  - name: alert.rules
    rules:
    - alert: HighCPUUsage
      expr: sum(rate(container_cpu_usage_seconds_total{image!=""}[2m])) by (pod) > 0.8
      for: 1m
      labels:
        severity: warning
      annotations:
        summary: "High CPU usage detected on pod {{ $labels.pod }}"
        description: "Pod {{ $labels.pod }} has CPU usage above 80% for more than 1 minute."
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please submit pull requests to improve the codebase or add new features.

