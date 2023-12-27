FROM golang:1.21.5 as builder
WORKDIR /app
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/bytez .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /bin/bytez /bin/bytez
CMD ["/bin/bytez"]
