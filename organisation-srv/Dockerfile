FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./organisation-srv /src/server/organisation-srv

RUN go get ./src/server/organisation-srv

ENTRYPOINT ["wait_container", "organisation-srv"]
