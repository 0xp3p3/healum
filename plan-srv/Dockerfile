FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./plan-srv /src/server/plan-srv

RUN go get ./src/server/plan-srv

ENTRYPOINT ["wait_container", "plan-srv"]
