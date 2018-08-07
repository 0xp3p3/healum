FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./task-srv /src/server/task-srv

RUN go get ./src/server/task-srv

ENTRYPOINT ["wait_container", "task-srv"]
