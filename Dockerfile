# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/go:latest as build

WORKDIR /work

COPY . .
RUN go vet -v
RUN go test -v -short

RUN CGO_ENABLED=0 go build -o envoy-control-plane

FROM cgr.dev/chainguard/static:latest

COPY --from=build /work/envoy-control-plane /envoy-control-plane
CMD ["/envoy-control-plane", "serve"]