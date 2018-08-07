FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv

RUN go get ./src/server/db-srv


ENTRYPOINT ["wait_container", "db-srv"]
