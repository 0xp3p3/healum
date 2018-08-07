FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./content-srv /src/server/content-srv

RUN go get ./src/server/content-srv

ENTRYPOINT ["wait_container", "content-srv"]
