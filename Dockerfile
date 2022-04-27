FROM golang:1 AS builder
WORKDIR /app
COPY . .
RUN make deps
RUN make build

FROM alpine:3 AS runner
WORKDIR /app
COPY --from=builder /app/bin/* app.bin
RUN chown -R nobody: /app
ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080
CMD ["./app.bin"]
