FROM golang:1.19.2

WORKDIR /home
COPY ./pkg /home

RUN cd /home && go build -o kubernetes

CMD ["/home/kubernetes"]