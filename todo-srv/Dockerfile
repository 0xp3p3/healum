FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./todo-srv /src/server/todo-srv

RUN go get ./src/server/todo-srv

ENTRYPOINT ["wait_container", "todo-srv"]
