FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./audit-srv /src/server/audit-srv

RUN go get ./src/server/audit-srv

ENTRYPOINT ["wait_container", "audit-srv"]
