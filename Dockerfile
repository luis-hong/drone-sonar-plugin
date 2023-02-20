FROM golang:1.20.0-alpine as build
RUN mkdir -p /go/src/github.com/aosapps/drone-sonar-plugin
WORKDIR /go/src/github.com/aosapps/drone-sonar-plugin 
COPY *.go ./
COPY api ./api/
COPY go.mod ./
COPY vendor ./vendor/
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o drone-sonar

FROM openjdk:11.0.8-jre

# https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.8.0.2856.zip
ARG SONAR_VERSION=4.8.0.2856
ARG SONAR_SCANNER_CLI=sonar-scanner-cli-${SONAR_VERSION}
ARG SONAR_SCANNER=sonar-scanner-${SONAR_VERSION}

RUN apt-get update \
    && apt-get install -y nodejs curl \
    && apt-get clean

COPY --from=build /go/src/github.com/aosapps/drone-sonar-plugin/drone-sonar /bin/
WORKDIR /bin

RUN curl https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/${SONAR_SCANNER_CLI}.zip -so /bin/${SONAR_SCANNER_CLI}.zip
RUN unzip ${SONAR_SCANNER_CLI}.zip \
    && rm ${SONAR_SCANNER_CLI}.zip 

ENV PATH $PATH:/bin/${SONAR_SCANNER}/bin

ENTRYPOINT /bin/drone-sonar
