FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./static-srv /src/server/static-srv

RUN go get ./src/server/static-srv

ENTRYPOINT ["wait_container", "static-srv"]
