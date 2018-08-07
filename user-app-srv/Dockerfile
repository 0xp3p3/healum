FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./user-app-srv /src/server/user-app-srv

RUN go get ./src/server/user-app-srv

ENTRYPOINT ["wait_container", "user-app-srv"]
