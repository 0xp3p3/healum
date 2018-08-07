FROM golang

ENV GOPATH /
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH

ADD ./db-srv /src/server/db-srv
ADD ./activity-srv /src/server/activity-srv
ADD ./common /src/server/common
ADD ./vendor /src/server/vendor
ADD ./plugins.go /src/server/

RUN chmod 777 ./src/server/plugins.go
RUN go get github.com/micro/micro

RUN go get ./src/server/activity-srv

RUN cp ./src/server/plugins.go ./src/github.com/micro/micro/
RUN go build -i -o $GOPATH/bin/micro $GOPATH/src/github.com/micro/micro/main.go $GOPATH/src/github.com/micro/micro/plugins.go

RUN chmod 777 ./src/server/activity-srv/wait_container
RUN cp ./src/server/activity-srv/wait_container ./bin/

ENTRYPOINT ["wait_container", "activity-srv"]
