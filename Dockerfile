ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd

COPY pkg pkg

RUN go build -o /terraform_ui ./cmd

FROM alpine:3.15 as app

EXPOSE 8080

# Add a user to run in non-root mode
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

USER nobody:nobody

COPY dist dist
COPY --from=builder /terraform_ui .

CMD [ "./terraform_ui", "server" ]
