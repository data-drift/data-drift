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
git clone https://github.com/data-drift/data-drift
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

5. Create a namespace for datadrift and apply the deployment resources using Kustomize:

```bash
kubectl create namespace datadrift
kubectl apply -k .
```

This will apply the following resources:

- Service: Defines how the DataDrift application will be exposed.
- Deployment: Deploys the DataDrift application containers to the Kubernetes cluster.
- Ingress: Configures traffic routing to the DataDrift application.

6.  Verify that the deployment was successful:

```bash
kubectl get pods -n datadrift
kubectl get services -n datadrift
kubectl get ingress -n datadrift
```

Make sure the pods are running, the service is exposed, and the ingress is properly configured.

7. Check the pod and the github credentials

Either go to the datadrift.yourdomain.com or bind a port to your localhost

```bash
kubectl port-forward datadrift-6df8ff84c5-bq2p5 8082:8080 -n datadrift
```

Go to your URL [localhost:8082](localhost:8082) you should see {"status":"OK"}.
Go to [localhost:8082/ghhealth](localhost:8082/ghhealth) you should see {"status":"OK"}.
