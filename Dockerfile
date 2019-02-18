FROM golang:latest

WORKDIR /go/src/soccer
COPY . .
RUN echo 'we are running some # of cool things'
RUN go get -d -v ./...
RUN go install -v ./...
#CMD [docker image rm $(docker image ls -a -q)]
CMD [docker image prune -a]
CMD ["soccer"]
