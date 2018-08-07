FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./response-srv /src/server/response-srv

RUN go get ./src/server/response-srv

ENTRYPOINT ["wait_container", "response-srv"]
