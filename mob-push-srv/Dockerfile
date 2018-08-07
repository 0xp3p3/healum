FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./mobpush-srv /src/server/mobpush-srv

RUN go get ./src/server/mobpush-srv

ENTRYPOINT ["wait_container", "mobpush-srv"]
