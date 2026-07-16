// Jenkins Pipeline untuk Portfolio Backend - CI/CD
// Build, Test, dan Deploy ke Docker

pipeline {
  agent any

  environment {
    PRODUCTION_SERVER_IP = "${env.PRODUCTION_SERVER_IP}"
    SONAR_HOST_URL       = 'https://sonar.bypur.my.id'
    DOCKER_IMAGE         = 'bypur-api-go'
    GOROOT               = '/usr/local/go'
    PATH                 = "${env.GOROOT}/bin:${env.PATH}"
    GOTOOLCHAIN          = 'go1.25.4'
  }

  stages {
    // 1. Get Source Code
    stage('Checkout') {
      steps {
        checkout scm
        echo 'Source code checked out'
      }
    }

    // 2. Go Build Check
    stage('Go Build Check') {
      steps {
        sh 'go version'
        sh 'go env'
        sh 'go build ./...'
        echo 'Go project compilation check passed'
      }
    }

    // 3. Run Tests and Generate Coverage
    stage('Go Test Coverage') {
      steps {
        sh 'go clean -testcache'
        sh 'go test ./... -coverprofile=coverage.out'
        echo 'Go test coverage report generated'
      }
    }

    // 4. SonarQube Analysis
    stage('SonarQube Analysis') {
      steps {
        script {
          withCredentials([string(credentialsId: 'sonarqube-token', variable: 'SONAR_TOKEN')]) {
            sh '''
              # Download SonarScanner CLI
              echo "Downloading SonarScanner CLI..."
              curl -sSLo sonar-scanner.zip https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-6.2.1.4610-linux-x64.zip
              
              # Extract the zip file
              echo "Extracting SonarScanner CLI..."
              unzip -o -q sonar-scanner.zip
              
              # Run the analysis
              echo "Running SonarQube Analysis..."
              ./sonar-scanner-6.2.1.4610-linux-x64/bin/sonar-scanner \
                -Dsonar.token="${SONAR_TOKEN}" \
                -Dsonar.host.url="${SONAR_HOST_URL}"
              
              # Clean up downloaded files
              echo "Cleaning up..."
              rm -rf sonar-scanner-6.2.1.4610-linux-x64 sonar-scanner.zip
            '''
          }
          echo 'SonarQube analysis completed'
        }
      }
    }

    // 5. Quality Gate Check
    stage('Quality Gate') {
      steps {
        script {
          timeout(time: 5, unit: 'MINUTES') {
            // Tunggu hasil quality gate dari SonarQube
            withCredentials([string(credentialsId: 'sonarqube-token', variable: 'SONAR_TOKEN')]) {
              def qg = sh(
                script: '''
                  STATUS="NONE"
                  for i in {1..12}; do
                    STATUS=$(curl -s -u "${SONAR_TOKEN}": "${SONAR_HOST_URL}/api/qualitygates/project_status?projectKey=bayupur-portofolio-be" | grep -o '"status":"[^"]*"' | head -n 1 | cut -d'"' -f4 || echo "NONE")
                    if [ "$STATUS" != "NONE" ] && [ ! -z "$STATUS" ]; then
                      break
                    fi
                    echo "Quality Gate status is still pending... waiting 5 seconds (attempt $i/12)..." >&2
                    sleep 5
                  done
                  echo "$STATUS"
                ''',
                returnStdout: true
              ).trim()
              
              echo "Quality Gate Status: ${qg}"
              
              if (qg != 'OK') {
                error "Quality Gate failed! Check SonarQube dashboard: ${SONAR_HOST_URL}/dashboard?id=bayupur-portofolio-be"
              }
              
              echo 'Quality Gate passed'
            }
          }
        }
      }
    }

    // 6. Build Docker Image
    stage('Build Docker Image') {
      steps {
        script {
          sh "docker build -t ${DOCKER_IMAGE}:latest ."
          echo "Docker image built successfully"
        }
      }
    }

    // 7. Deploy to Docker Container
    stage('Deploy') {
      steps {
        script {
          withCredentials([
            string(credentialsId: 'database-url', variable: 'DATABASE_URL'),
            string(credentialsId: 'jwt-secret', variable: 'JWT_SECRET'),
            string(credentialsId: 'cors-origin', variable: 'CORS_ORIGIN'),
            string(credentialsId: 'x-api-key', variable: 'X_API_KEY')
          ]) {
            sh '''
              # Stop dan hapus container lama jika ada
              docker stop bypur_api_go || true
              docker rm bypur_api_go || true
              
              # Jalankan container baru terhubung ke bypur_network
              docker run -d \
                --name bypur_api_go \
                --network bypur_network \
                --memory="512m" \
                --cpus="0.5" \
                --restart=unless-stopped \
                -p 3001:3001 \
                -e APP_ENV=production \
                -e SERVER_PORT=3001 \
                -e DB_URL="${DATABASE_URL}" \
                -e JWT_SECRET="${JWT_SECRET}" \
                -e SERVER_CORS_ORIGIN="${CORS_ORIGIN}" \
                -e SECURITY_X_API_KEY="${X_API_KEY}" \
                bypur-api-go:latest
              
              # Tunggu container start
              sleep 5
              
              # Health check
              curl -f http://localhost:3001/health || exit 1
              
              echo "Deployment backend berhasil! Container running on port 3001"
            '''
          }
        }
      }
    }
  }

  // Post-build Actions
  post {
    always {
      // Bersihkan workspace menggunakan deleteDir() bawaan Jenkins Pipeline
      deleteDir()
    }
    success {
      echo 'Pipeline completed successfully!'
    }
    failure {
      echo 'Pipeline failed! Check the logs above.'
    }
  }
}
