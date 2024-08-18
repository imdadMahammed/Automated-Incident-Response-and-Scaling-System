pipeline {
    agent any
    stages {
        stage('Build Docker Image') {
            steps {
                script {
                    docker.build('sample-app:latest')
                }
            }
        }
        stage('Push to Docker Hub') {
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'dockerhub-credentials') {
                        docker.image('sample-app:latest').push()
                    }
                }
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                kubernetesDeploy(
                    kubeconfigId: 'kubeconfig-credentials',
                    configs: 'k8s/*.yaml'
                )
            }
        }
    }
}
