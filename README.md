# Prequisites for Multi-Master Kubernetes Cluster setup:
1. 2 GB or more of RAM per machine (any less will leave little room for your apps)
2. 2 CPUs or more
3. Full network connectivity between all machines in the cluster (public or private network is fine)
4. Unique hostname, MAC address, and product_uuid for every node. This may arise when cluster is installed on VM's.
5. Certain ports are open on your machines. Ensure the port required are unused.
6. Swap disabled. You MUST disable swap in order for the kubelet to work properly
7. Install any load balancer to route the traffic from worker nodes to api-servers in master cluster equally. Capture the Load Balancer IP(This will be useful during kubeadm init).
8. Below setup instructions were tried and tested on Ubuntu 16.04, should work for all debian linux flavours. 
9. Cluster is setup in GCE.

# HA Kubernetes Cluster setup instructions:

1. Install docker (https://docs.docker.com/install/linux/docker-ce/ubuntu/)
2. Install kubelet kubeadm and kubelet (https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
3. Copy below into kubeadm-config.yaml file
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: stable
#REPLACE with `loadbalancer` IP
controlPlaneEndpoint: "10.128.0.8:6443"  // If loadbalancer is deployed then provide load balancer IP here.
networking:
  podSubnet: 192.168.0.0/18
4. Run below command to bring up kubernetes cluster which includes core-dns, kube-proxy, etcd, kube-apiserver, kube-scheduler, kube-controller-manager, kubelet.
#kubeadm init --config=/etc/kubernetes/kubeadm/kubeadm-config.yaml --upload-certs
5. Above command will provides two commands as output. One with --control-plane used by other nodes to join as master and another without --control-plane will be used by other nodes to join as worker nodes.
Note: Tokens will expire by default 1 hour for security purpose. By passing --token-ttl 0, we can make the token valid indefinite time period.(Not suggested)
6. Execute below commands to run kubectl commands on master node:
  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config
7. kubectl commands can be run from other nodes by copying above config file and by following about steps.
8. Choose the network plugin over which pods communicate inside the cluster. For calico,
#kubectl apply -f https://gist.githubusercontent.com/joshrosso/ed1f5ea5a2f47d86f536e9eee3f1a2c2/raw/dfd95b9230fb3f75543706f3a95989964f36b154/calico-3.5.yaml
9. Your HA Kubernetes cluster must be up and running.
10. To Validate HA by bringing down one the master node and deploy any sample application. 

# Create a development namespace
 command: kubectl create ns development
 
# Install and configure Helm in Kubernetes:
Helm Client Installation on Jenkins server and Tiller(Helm Server) on Kubernetes Cluster:
1. Download the latest version of Helm from https://github.com/kubernetes/helm/releases
2. Unpack the archive:
$ gunzip helm-v2.8.1-linux-amd64.tar.gz
$ tar -xvf helm-v2.8.1-linux-amd64.tar
$ sudo mv l*/helm /usr/local/bin/.
3. Run below command install helm components
#helm init
4. Run below commands to setup kubernetes credentials:
#kubectl create serviceaccount --namespace kube-system tiller
#kubectl create clusterrolebinding tiller-cluster-rule \
   --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
#kubectl patch deploy --namespace kube-system tiller-deploy \
   -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

# Jenkins setup:

1. Follow the link to setup Jenkins on GCE, https://linuxize.com/post/how-to-install-jenkins-on-debian-9/
2. Create CI/CD Pipeline by providing git repo https://github.com/laxmanvallandas/examples.git. Enable the option to "GitHub hook trigger for GITScm polling", to start executing the [Jenkinsfile](Jenkinsfile) upon receiving updates on the repo.
3. Also add the web-hook in github project settings, 
http://<JenkinsServerIP>:8080/github-webhook/
  Github sends a callback to Jenkins via this URL to trigger the pipeline corresponding to it.

# Use Helm to deploy the application on Kubernetes Cluster from CI server.
1. Jenkinsfile is the MAIN AUTOMATED SCRIPT written in Groovy, to pipeline the events starting from 
  docker build - to build the image
  docker push - to push to private docker repository
  Using helm, Deploy Guestbook application in kubernetes cluster inside development namespace.
  Using helm, Upgrade Guestbook Application.
  Rollback of Application is not automated, as part of this script, but can be triggered manually upon failure.
  Sample command below:
  #helm del --purge guestbook --tiller-namespace development
 
 # Create a monitoring namespace in the cluster.
command: kubectl create ns monitoring

# Setup Prometheus using Helm:
helm install stable/prometheus --name prom --values guestbook/prometheus/prom.values --namespace monitoring --tiller-namespace development

# Setup Grafana using Helm:
helm install stable/grafana --tiller-namespace development --namespace monitoring --values guestbook/prometheus/grafana.values --name grafana

# Add Prometheus as Data Source in Grafana Dashboard to get the statistics of kubernetes cluster in the UI.
1. Add Prometheus IP address in the URL and add to Grafana Dashboard.
2. Search for Kubernetes deployment daemonset metrics, copy the code from grafana website and import it in Grafana.

# Finally, modify guestbook/php-redis/index.html with below change to set the background color to guestbook application.
<style>
body {
background-color: #d24dff
}
</style>

  commit the change
  
  REFRESH guestbook URL to Observe the change. 



Note: All files and modifications required to deploy guestbook application has been checked into 
https://github.com/laxmanvallandas/examples.git
