FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./cloudkey-srv /src/server/cloudkey-srv

RUN go get ./src/server/cloudkey-srv

ENTRYPOINT ["wait_container", "cloudkey-srv"]
