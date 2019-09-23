pipeline {
    agent {
        docker {
            image 'agilesolutions/bomverifier'
        }
    }
    stages {
        stage('Build') {
            steps {
                sh 'bomverifier https://raw.githubusercontent.com/agilesolutions/bomverifier/master/bom.yaml'
            }
        }
    }
}