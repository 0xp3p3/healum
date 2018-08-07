FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./team-srv /src/server/team-srv

RUN go get ./src/server/team-srv

ENTRYPOINT ["wait_container", "team-srv"]
