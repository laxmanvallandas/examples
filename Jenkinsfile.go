node{
  //Define all variables
  def project = 'laxman'
  def appName = 'guestbook'
  def serviceName = "${appName}-backend"  
  def imageVersion = 'development'
  def namespace = 'development'
  def imageTag = "${appName}:${imageVersion}.${env.BUILD_NUMBER}"
  
  //Checkout Code from Git
  checkout scm
  
  //Stage 1 : Build the docker image.
  stage('Build image') {
    sh("docker build -t ${project}/${imageTag} -f /var/lib/jenkins/workspace/guestbook/guestbook/php-redis/Dockerfile .")
  }
  
  //Stage 2 : Push the image to docker registry
  stage('Push image to registry') {
      sh("docker login -u laxman -p ${DOCKER_HUB}")
      sh("docker push ${project}/${imageTag}")
  }
  
  //Stage 3 : Deploy Application
  stage('Deploy Application') {
       switch (namespace) {
              //Roll out to Dev Environment
              case "development":
                   // Create namespace if it doesn't exist
                   sh("sed -i.bak 's/IMAGE-TAG/${imageTag}/g' guestbook/php-redis/helm-chart/templates/guestbook-all-in-one.yaml")
                   sh("helm install --name guestbook --tiller-namespace development --namespace development ./guestbook/php-redis/helm-chart")
                   //sh("echo http://`kubectl --namespace=${namespace} get service/guestbook --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > guestbook")
                   break
  }
}
}
