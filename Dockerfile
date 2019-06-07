FROM golang:alpine

EXPOSE 80

ADD ./ /tasksREST
WORKDIR /tasksREST


RUN apk add git
RUN go mod download
RUN go install -v ./cmd/tasks

ENTRYPOINT ["tasks"]


