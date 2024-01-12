ARG GO_VERSION=1.21.6

# STAGE 1
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src/

COPY go.mod ./
RUN go mod download

COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/psych cmd/psych/main.go

# STAGE 2
FROM gcr.io/distroless/static-debian11:nonroot

LABEL maintainer="brittonhayes"
LABEL org.opencontainers.image.source="https://github.com/brittonhayes/psych"
LABEL org.opencontainers.image.description="Find a mental health professional."
LABEL org.opencontainers.image.licenses="MIT"

COPY --from=builder --chown=nonroot:nonroot /bin/psych /bin/psych

EXPOSE 8080

ENTRYPOINT [ "/bin/psych" ]

CMD ["/bin/psych"]