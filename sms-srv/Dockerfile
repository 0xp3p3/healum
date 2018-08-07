FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./sms-srv /src/server/sms-srv

RUN go get ./src/server/sms-srv

ENTRYPOINT ["wait_container", "sms-srv"]
