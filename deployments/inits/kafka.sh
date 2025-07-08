#! /bin/sh

kafka-topics --create --if-not-exists \
  --topic gate \
  --partitions 10 \
  --replication-factor 1 \
  --bootstrap-server kafka:9092

kafka-topics --create --if-not-exists \
  --topic log \
  --partitions 10 \
  --replication-factor 1 \
  --bootstrap-server kafka:9092