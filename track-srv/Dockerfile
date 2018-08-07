FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./track-srv /src/server/track-srv

RUN go get ./src/server/track-srv

ENTRYPOINT ["wait_container", "track-srv"]
