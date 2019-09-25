# Kubernetes HELM configuration injector
Fetches secure Spring Boot aplication configs from Spring Cloud Config server and injects config on HELM chart on a k8s configmap. This GO app get packaged on docker images and executed as Jenkins Docker Agent on Jenkins Groovy deployment pipeline
[Read about Spring Cloud Config server](https://o7planning.org/en/11727/understanding-spring-cloud-config-client-with-example)

## functionality

1. wget the BOM txt file from github
2. go into the springboot jar file zip file and discover all libraries
3. check the compliancy againt the BOM txt
4. report and conditionally break off

## setup

* [goto](https://www.katacoda.com/courses/docker/deploying-first-container)
* git clone https://github.com/agilesolutions/configinjector.git
* curl -LO https://dl.google.com/go/go1.13.linux-amd64.tar.gz
* tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
* export PATH=$PATH:/usr/local/go/bin	
* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH

## build

```
go build -o  .

configinjector -url=https://github.com/o7planning/spring-cloud-config-git-repo-example -directory=chart

docker build -t agilesolutions/configinjectorr:latest .
```

## run
configinjector -url=https://github.com/o7planning/spring-cloud-config-git-repo-example -directory=chart

## where to find the Springboot BOM details and release trains

* [spring-boot-dependencies BOM](https://github.com/spring-projects/spring-boot/blob/master/spring-boot-project/spring-boot-dependencies/pom.xml)

## now run this docker agent on a jenkins pipeline, lets spin up jenkins

* [go to katacoda](https://www.katacoda.com/courses/kubernetes/helm-package-manager)
* create directory /jenkins
* docker run -d --name jenkins --user root --privileged=true -p 8080:8080 -v /jenkins:/var/jenkins_home -v /var/run/docker.sock:/var/run/docker.sock jenkinsci/blueocean
* docker logs -f jenkins
* docker exec -ti jenkins bash
* docker ps -a
* browse to http://localhost:8080 and wait until the Unlock Jenkins page appears.
* get password from /jenkins/secrets/initialAdminPassword
* create new pipeline job from https://github.com/agilesolutions/bomverifier.git

## include on pipeline

```
pipeline {
  agent none
  environment {
    DOCKER_IMAGE = null
  }
  stages {
    stage('config') {
      agent {
          docker {
              image 'agilesolutions/configinjector:latest'
          }
      }
      steps {
        sh 'bomverifier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate'
      }
    }
    stage("deploy") {
      when {
        branch "master"
      }
      steps {
        container("helm") {
			k8sUpgrade(params.artifact, params.version)
        }
      }
    }
```
Groovy k8sUpgrade pipeline step...

```
def call(artifact, version ) {
    sh 'helm upgrade \
        ${artifact} \
        charts/${artifact} -i \
        --namespace ${artifact} \
        --set image.tag=${version} \
        --set ingress.host=${artifact} \
        --reuse-values'
}
```

## read

1 [check this](https://www.callicoder.com/docker-golang-image-container-example/)
2 [parse yaml](https://stackoverflow.com/questions/28682439/go-parse-yaml-file/28683173)
3 [wget to file](https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go)
4 [get go package](https://gopkg.in/yaml.v2)
5 [jenkins pipelines and docker agents](https://jenkins.io/doc/book/pipeline/docker/)
6 []()
