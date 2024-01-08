#스테이지 1 : go 환경에서 바이너리 컴파일
FROM golang:1.21.5 AS build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs

#스테이지 2 : 이미지 빌드
#scratch 이미지 : 아무런 배포파일도 포함하지 않는 이미지
FROM scratch 
COPY --from=build /src/kvs .
EXPOSE 8080
CMD ["/kvs"]