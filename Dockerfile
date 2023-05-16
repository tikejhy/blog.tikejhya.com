# syntax=docker/dockerfile:1

FROM golang:1.19 AS BuildStage

WORKDIR /app

COPY go.mod ./
RUN go mod download
RUN go get -u github.com/gin-gonic/gin

COPY *.go ./

EXPOSE 8080

# CGO has to be disabled for alpine
RUN CGO_ENABLED=0 go build -o /blog

#docker

FROM alpine:latest

WORKDIR /

COPY . .
COPY --from=BuildStage /blog /blog

EXPOSE 8080

CMD [ "/blog" ]
