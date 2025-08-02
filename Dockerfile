FROM golang:1.24-bookworm AS builder
WORKDIR /app/cookbook

COPY . .

RUN CGO_ENABLED=0 go build -o cookbook
RUN echo "Built cookbook"

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y curl

COPY --from=builder /app/cookbook/cookbook /usr/local/bin/
COPY --from=builder /app/cookbook/db/migrate /opt/cookbook/db/migrate

ENV USER=cookbook
ENV UID=1000
ENV GID=1000
RUN addgroup --gid "$GID" "$USER" \
  && adduser \
  --disabled-password \
  --gecos "cookbook" \
  --home "/opt/$USER" \
  --ingroup "$USER" \
  --no-create-home \
  --uid "$UID" \
  "$USER" \
  && chown "$USER" /usr/local/bin/cookbook \
  && chmod u+x /usr/local/bin/cookbook

WORKDIR "/opt/$USER"
RUN chown cookbook "/opt/$USER"
USER cookbook

ENTRYPOINT ["/usr/local/bin/cookbook"]
CMD ["run"]
