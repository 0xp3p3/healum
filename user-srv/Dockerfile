FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./user-srv /src/server/user-srv

RUN go get ./src/server/user-srv

ENTRYPOINT ["wait_container", "user-srv"]
