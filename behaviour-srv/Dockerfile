FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./behaviour-srv /src/server/behaviour-srv

RUN go get ./src/server/behaviour-srv

ENTRYPOINT ["wait_container", "behaviour-srv"]
