##  Asignment 2
##### 1. Create  a Kubernetes cluster which runs a GoLang service and reponds with the message "Hello world" on accesing the API.

Procedure:
1. Implement the script to write the message "HelloWorld".
2. Generate a Docker image and upload it to the local repository.
3. Create a deployment in Kubernetes which uses the generated image.
4. Expose the deployment to access the endpoint.
5. Deploy a debug-pod(OS container) to validate the functionality of the HelloWorld container within the cluster.

Alternatively, the pod's exposed port may be port-forwarded to be reachable from outside the container for validation.

Commands to get the task done:

##### Generate the image to be deployed in Kubernetes cluster.

```bash
eval $(minikube docker-env) # to export the .bash variables to re-use the same in the minikube enviroment.

docker build --tag=localhost:5000/helloworldgo .
docker run -d --publish=8080:8080 localhost:5000/helloworldgo:latest
docker push localhost:5000/helloworldgo:latest`
```

##### Commands to deploy the generated image on a Kubernetes cluster

```bash
kubectl create deployment hello-world --image=localhost:5000/helloworldgo:latest --replicas=1
kubectl scale deployment hello-world --replicas=2
kubectl expose deployment hello-world --port=8080
kubectl port-forward <POD-NAME> PORT:PORT
```
Validation environment:
* minikube : v1.25.2
* docker : 20.10.12
* kubectl :v1.24.0
* Ubuntu : 22.04

**Deploy the pods and expose the deployment as a service.:**
[![Pod Deployment](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/createdeployment%20and%20svc.png?raw=true "Pod Deployment")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/createdeployment%20and%20svc.png?raw=true "Pod Deployment")

**Validate the functionality of the HelloWorld service from Ubuntu pod using cURL:**
[![TestingHelloWorldFromUbuntuPod](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/internalcurl.png?raw=true "TestingHelloWorldFromUbuntuPod")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/internalcurl.png?raw=true "TestingHelloWorldFromUbuntuPod")

[![LogsAtPod](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/ServerLogs.png?raw=true "LogsAtPod")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/ServerLogs.png?raw=true "LogsAtPod")

** Expose the pod by port-forwarding and access the /helloworld endpoint service from the VM that hosts the cluster**

[![ExposePort](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/pods+portfwd.png?raw=true "ExposePort")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/pods+portfwd.png?raw=true "ExposePort")

[![CurlFromVM](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/curltofwdport.png?raw=true "CurlFromVM")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/curltofwdport.png?raw=true "CurlFromVM")


#### Sub-questions:

##### 1. How do you make the service Scalable?

A service can be made scalable by accounting for the following factors:
1. Avoiding single point of failures for services.
2. Distribute the core logic into multiple instances.
3. Incorporate statelessness in the introduced service, which enables easy recovery in case of failures.
4. Caching data reduces the time taken for both fetching and computing.
5. The resources for the pods deployed must be sufficient and effictively used, to ensure maximum utilization of the provisioned resources and minimum provisioning of resources.

The instances/services in Kubernetes can be either manually or autoscaled by the HPA  based on the CPU utilization across instances. This ensures that the requests are served at all times.

kubectl autoscale deployment [deployment_name ]--cpu-percent=50 --min=1 --max=10
The parameter cpu-percent determines the load on the CPU beyond which the HPA creates a new pod to handle requests.


##### 2. What CI/CD pipeline would you use?(If not code, please describe every step from commit of new code until the new code is running in production)

The steps involved from CI to CD are as follows:
1. The developer branches from master, either to provide a fix for a defect, or to add a new feature. During this process, a corresponding branch is created and checkout.
2. Once the feature is ready to be merged, a jenkins git-hook which is already set triggers a build on the Jenkins server to validate that the existing features are not broken and the new feature/bugfix works along with the existing changes. This can be validated from the Unit Tests provided by the developer.
3. If all the unit tests have passed successfully, the build can now be pushed to artifactory where automations, if present push the same onto a QA environment where further rounds of testing are done, with both infrastructural and feature/functional level tests.
4. Once the features are tested locally, the changes are then pushed to master with the retest of the feature/bugfix revalidated in the staging environment.
5. Post successful validation, the master-build with all features bundled and built, is then moved to the CD department, where the builds are automated to be uploaded to the production environment.

##### 3.How would you store and deploy secrets(Such as API Keys)

Kubernetes has the provision to provide sensitive data such as tokens and keys which may be required for the purpose of authentication across the services. Kuberenetes secrets are similar to config-maps, whereas the fields in the case of secrets are encoded and stored.

The sensitive information is encoded in Base64 and can be injected into the environment of pods. 

A secret can be handled using the following commands:
```bash
kubectl get secret
kubectl get secret <Secret_Name> -o jsonpath='{.data.*}' | base64 -d
kubectl create secret generic <secret_name> --from-literal=TEST_USER=secretvalue
kubectl delete secret <secret_name>
```

There are multiple ways by which a secret can be injected and used in the Kubernetes cluster.
1. As a file in a volume mounted to the the containers.
2. As a container environment variable
3. By the kubelet while pulling images for the pod

The case of adding secrets as an environment variable has been experimented.

**Define the secrets that need to be used using kubectl create secret**
[![creatingSecrets](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/NewSecretAddition.png?raw=true "creatingSecrets")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/NewSecretAddition.png?raw=true "creatingSecrets")

**Displaying the secret which is base64 encoded**
[![ShowSecrets](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/SecretDisplay.png?raw=true "ShowSecrets")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/SecretDisplay.png?raw=true "ShowSecrets")

**Exposing the secret to a deployment**
[![ExposingSecrets](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/loadedfromSecretsinDeployment.png?raw=true "ExposingSecrets")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/loadedfromSecretsinDeployment.png?raw=true "ExposingSecrets")

**Ensuring that the env variables in the Pod's environment are injected -SECRET_USERNAME and SECRET_PASSWORD, derived from USERNAME and PASSWORD of *testsecret***
[![ExposedEnvVars](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/injectedintoPodenv.png?raw=true "ExposedEnvVars")](https://github.com/kishenv/AssignmentTwo/blob/main/Screenshots/injectedintoPodenv.png?raw=true "ExposedEnvVars")


##### 4. How do you test how well your infrastructure scales (When many requests come in)?

The tests associated with the infrastructure can be load, endurance or stress tests.

In the case of a web-application, a service can be deployed internally to spam requests to the API server, this would inturn cause the load to increase, and ultimately scale the deployments in the cluster, provided an autoscale policy is defined with the HPA.

eg: kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done" would generate requests to the php-endpoint. This would cause an increase in the CPU usage, once the load is beyond the limits set while defining the HPA, the required number of replicas are spawned to serve the requests.

##### 5. How do you provide an SSL certificate for your service?

The SSL certificate can either be bundled along with the docker image during build. This would load the file in the appropropriate directory. The command update-ca-certificates which is run as one of the steps in Dockerfile appends the certificate related information to /etc/ssl/certs/ca-certificates.crt.

Following are the steps to be added to the dockerfile to add SSL certificate to a Pod:

```bash
ADD your_ca_root.pem /usr/local/share/ca-certificates/foo.pem
RUN chmod 644 /usr/local/share/ca-certificates/foo.pem 
RUN update-ca-certificates
```

An alternate way to use certificates is by creating a secret/config-map with TLS and creating a Volume and a VolumeMount in order to load the certificate in service.

The following can be done to achieve the same:

1. To be added under pod-level when the certificate is loaded from Config-map:
```yaml
         volumeMounts:
         - name: ca-pemstore
         mountPath: /etc/ssl/certs/my-cert.pem
          subPath: my-cert.pem
          readOnly: false
      volumes:
      - name: ca-pemstore
        configMap:
          name: ca-pemstore
```

2. As Secret:

```yaml
        volumeMounts:
          - mountPath: "/etc/ssl/certs"
            name: test-ssl
            readOnly: true
      volumes:
        - name: test-ssl
          secret:
            secretName: test-ssl
```

Generating a key+pem locally which can be ported to the cluster.

```bash
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes  -days 365 key ca.key -out ca.crt -subj "/CN=yourdomain.com"
kubectl create secret tls my-tls-secret \
--key < private key filename> \
--cert < certificate filename>
```


