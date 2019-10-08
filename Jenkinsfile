// vim: set filetype=groovy:


def buildName = "${env.JOB_BASE_NAME.replaceAll("%2F", "-").replaceAll("\\.", "-").take(20)}-${env.BUILD_ID}"

pipeline {
    agent {
      kubernetes {
        cloud 'zeebe-ci'
        label "zeebe-ci-build_${buildName}"
        defaultContainer 'jnlp'
        yamlFile '.ci/podSpecs/distribution.yml'
      }
    }

    environment {
      NEXUS = credentials("camunda-nexus")
      SONARCLOUD_TOKEN = credentials('zeebe-sonarcloud-token')
    }

    options {
        buildDiscarder(logRotator(daysToKeepStr: '-1', numToKeepStr: '10'))
        timestamps()
        timeout(time: 45, unit: 'MINUTES')
    }

    stages {
        stage('Prepare') {
            steps {
                container('maven') {
                    sh '.ci/scripts/distribution/prepare.sh'
                }
            }
        }

        stage('Build (Java)') {
            steps {
                container('maven') {
                    sh '.ci/scripts/distribution/build-java.sh'
                }
            }
        }

        stage('Test (Java)') {
            parallel {
                stage('Analyse (Java)') {
                      steps {
                          container('maven') {
                              sh '.ci/scripts/distribution/analyse-java.sh'
                          }
                      }
                }

                stage('Unit (Java)') {
                    steps {
                        container('maven') {
                            sh '.ci/scripts/distribution/test-java.sh'
                        }
                    }
                }

                stage('IT (Java)') {
                    steps {
                        container('maven') {
                            sh '.ci/scripts/distribution/it-java.sh'
                        }
                    }
                }
            }

            post {
                always {
                    junit testResults: "**/*/TEST-*.xml", keepLongStdio: true
                }
            }
        }

        stage('Upload') {
            when { branch 'develop' }
            steps {
                container('maven') {
                    sh '.ci/scripts/distribution/upload.sh'
                }
            }
        }

        stage('Post') {
            parallel {
                stage('Docker') {
                    when { branch 'develop' }

                    environment {
                        VERSION = readMavenPom(file: 'parent/pom.xml').getVersion()
                    }

                    steps {
                        build job: 'zeebe-docker', parameters: [
                            string(name: 'BRANCH', value: env.BRANCH_NAME),
                            string(name: 'VERSION', value: env.VERSION),
                            booleanParam(name: 'IS_LATEST', value: env.BRANCH_NAME == 'master')
                        ]
                    }
                }

                stage('Docs') {
                    when { anyOf { branch 'master'; branch 'develop' } }
                    steps {
                        build job: 'zeebe-docs', parameters: [
                            string(name: 'BRANCH', value: env.BRANCH_NAME),
                            booleanParam(name: 'LIVE', value: env.BRANCH_NAME == 'master')
                        ]
                    }
                }
            }
        }
    }

    post {
        always {
            // Retrigger the build if there were connection issues
            script {
                build job: currentBuild.projectName, propagate: false, quietPeriod: 10, wait: false
            }
        }
    }
}

boolean connectionProblem() {
  return currentBuild.rawBuild.getLog(500).join('') ==~ /.*(ChannelClosedException|KubernetesClientException|ClosedChannelException|Connection reset).*/
}
