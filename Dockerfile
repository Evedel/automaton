FROM golang:1.21.1 AS build
COPY . /app
WORKDIR /app
RUN go build -o actioneer cmd/main.go
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl

FROM debian:bookworm-slim AS k8s
COPY --from=build /app/actioneer /app/actioneer
COPY --from=build /usr/local/bin/kubectl /usr/local/bin/kubectl
ENTRYPOINT [ "/app/actioneer" ]
