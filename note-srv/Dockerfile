FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./note-srv /src/server/note-srv

RUN go get ./src/server/note-srv

ENTRYPOINT ["wait_container", "note-srv"]
