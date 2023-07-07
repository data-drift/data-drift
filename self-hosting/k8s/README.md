# DataDrift Project

This project contains the necessary files to deploy the DataDrift application using Kustomize.

## Prerequisites

Make sure you have the following tools installed before getting started:

- Kustomize
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- An available Kubernetes cluster

## Deployment

To deploy the DataDrift application, follow the steps below:

1.  Clone this repository to your local machine:

```bash
git clone https://github.com/your-username/datadrift.git
```

2.  Navigate to the project directory:

```bash
cd data-drift/self-hosting/k8s
```

3.  Configure your secrets by creating a `.env.secret` file. You can use the provided `.env.secret.example` template as a starting point:

```bash
cp .env.secret.example .env.secret
```

Modify the `.env.secret` file by adding the appropriate values for your secrets.

4. Configure the ingress

In the ingress file, replace the `-host: datadrift.REPLACE_WITH_YOUR_DOMAIN` in the rules section
In the ingress file, replace the `-datadrift.REPLACE_WITH_YOUR_DOMAIN` in hosts section

5. Apply the deployment resources using Kustomize:

```bash
kubectl apply -k .
```

This will apply the following resources:

- Service: Defines how the DataDrift application will be exposed.
- Deployment: Deploys the DataDrift application containers to the Kubernetes cluster.
- Ingress: Configures traffic routing to the DataDrift application.

6.  Verify that the deployment was successful:

```bash
kubectl get pods
kubectl get services
kubectl get ingress
```

Make sure the pods are running, the service is exposed, and the ingress is properly configured.

## Contributions

Contributions to this project are welcome. If you would like to make improvements or report issues, please submit a pull request or open an issue.
