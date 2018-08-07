FROM healum/server:v2

ADD ./api /src/server/api
ADD ./user-srv /src/server/user-srv
ADD ./db-srv /src/server/db-srv
ADD ./cloudkey-srv /src/server/cloudkey-srv
ADD ./organisation-srv /src/server/organisation-srv

ADD ./common /src/server/common


RUN go get ./src/server/api

ENTRYPOINT ["wait_container", "api"]
