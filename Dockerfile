FROM golang
RUN go get -u google.golang.org/grpc
RUN useradd -ms /bin/bash gouser
RUN mkdir /golang-grpc
WORKDIR /golang-grpc
ADD /main/ /golang-grpc/
RUN go build -o main .
USER gouser
EXPOSE 8080
CMD ["./main"]
