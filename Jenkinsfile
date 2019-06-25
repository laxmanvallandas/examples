node{

  //Define all variables
  def project = 'laxmanvallandas'
  def appName = 'guestbook-go'
  def serviceName = "${appName}-backend"  
  def imageVersion = 'development'
  def namespace = 'development'
  def imageTag = "${project}/${appName}:${imageVersion}.${env.BUILD_NUMBER}"
  
  //Checkout Code from Git
  checkout scm
  
  //Stage 1 : Build the docker image.
  stage('Build image') {
      sh("cd guestbook-go")
      sh("docker build -t ${imageTag} .")
  }
  
  //Stage 2 : Push the image to docker registry
  stage('Push image to registry') {
      sh("docker login -u laxman -p ${DOCKER_HUB}")
      sh("docker push ${imageTag}")
  }
  
  //Stage 3 : Deploy Application
  stage('Deploy Application') {
       switch (namespace) {
              //Roll out to Dev Environment
              case "development":
                   // Create namespace if it doesn't exist
                   sh("sed -i.bak 's#/${project}/${appName}:${imageVersion}#${imageTag}#' guestbook-go/helm-chart/templates/guestbook-controller.json")
                   sh("helm install --name guestbook ./guestbook-go/helm-chart")
                   sh("echo http://`kubectl --namespace=${namespace} get service/${feSvcName} --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > ${feSvcName}")
                   break
  }

}
}
