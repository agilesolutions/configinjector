# BOM Bill of Material verifier
Scan Spring boot jar file for libraries complying to the content of bom.txt. This app is going to be wrapped on container and be run as a jenkins pipeline 2.0 agent.
Let the jenkins build fail if any of the included libraries on that spring boot app are violating the compliancy of the BOM test.
## functionality

1. wget the BOM txt file from github
2. go into the springboot jar file zip file and discover all libraries
3. check the compliancy againt the BOM txt
4. report and conditionally break off

## setup

* [goto](https://www.katacoda.com/courses/docker/deploying-first-container)
* git clone https://github.com/agilesolutions/bomverifier.git
* curl -LO https://dl.google.com/go/go1.13.linux-amd64.tar.gz
* tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
* export PATH=$PATH:/usr/local/go/bin	
* export GOPATH=/root/go
* export GOBIN=/usr/local/go/bin
* export PATH=$PATH:$(go env GOPATH)/bin
* go env GOPATH

## build

```
go build -o bomverifier .

bomverifier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate

docker build -t agilesolutions/bomverifier:latest .
```

## run
bomverfier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate

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
    stage('Verify') {
      agent {
          docker {
              image 'agilesolutions/bomverifier:latest'
          }
      }
      steps {
        sh 'bomverifier -url=https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.txt -terminate'
      }
    }
    stage('Build') {
      agent {
          docker {
              image 'maven:3-alpine'
            // do some caching on maven here
              args '-v $HOME/.m2:/root/.m2'
          }
      }
      steps {
        sh 'mvn clean install'
      }
    }
    stage('dockerbuild') {
      steps {
        script {
          DOCKER_IMAGE = docker.build("katacodarob/demo:latest")
        }
      }
    }
```


## read

1 [check this](https://www.callicoder.com/docker-golang-image-container-example/)
2 [parse yaml](https://stackoverflow.com/questions/28682439/go-parse-yaml-file/28683173)
3 [wget to file](https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go)
4 [get go package](https://gopkg.in/yaml.v2)
5 [jenkins pipelines and docker agents](https://jenkins.io/doc/book/pipeline/docker/)
6 []()
