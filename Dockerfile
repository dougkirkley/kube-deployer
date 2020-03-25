FROM golang as builder

LABEL maintainer="Douglass Kirkley <doug.kirkley@gmail.com"

# Create appuser.
ENV USER=kube
ENV UID=10001 

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"


WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# Build the Go app
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/kube-deployer

FROM alpine

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/kube-deployer /go/bin/kube-deployer

ENV VERSION="3.0.3"
# ENV BASE_URL="https://storage.googleapis.com/kubernetes-helm"
ENV BASE_URL="https://get.helm.sh"
ENV TAR_FILE="helm-v${VERSION}-linux-amd64.tar.gz"

RUN apk add --update --no-cache curl ca-certificates && \
    curl -L ${BASE_URL}/${TAR_FILE} |tar xvz && \
    mv linux-amd64/helm /usr/bin/helm && \
    chmod +x /usr/bin/helm && \
    rm -rf linux-amd64 && \
    apk del curl && \
    rm -f /var/cache/apk/*

# Expose port 9090 to the outside world
EXPOSE 9090
 
USER kube:kube

# Command to run the executable
ENTRYPOINT ["/go/bin/kube-deployer"]
