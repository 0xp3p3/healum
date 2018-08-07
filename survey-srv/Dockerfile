FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./survey-srv /src/server/survey-srv

RUN go get ./src/server/survey-srv

ENTRYPOINT ["wait_container", "survey-srv"]
