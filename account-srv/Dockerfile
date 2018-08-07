FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./account-srv /src/server/account-srv

RUN go get ./src/server/account-srv

ENTRYPOINT ["wait_container", "account-srv"]
