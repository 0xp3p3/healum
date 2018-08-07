FROM healum/server:v2

ADD ./common /src/server/common

ADD ./db-srv /src/server/db-srv
ADD ./product-srv /src/server/product-srv

RUN go get ./src/server/product-srv

ENTRYPOINT ["wait_container", "product-srv"]
