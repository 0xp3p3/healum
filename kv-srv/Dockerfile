FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./kv-srv /src/server/kv-srv

RUN go get ./src/server/kv-srv

ENTRYPOINT ["wait_container", "kv-srv"]
