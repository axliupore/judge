#!/bin/bash

HTTP="6048"
NSQ=""
WS=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --http=*)
      HTTP="${1#*=}"
      shift
      ;;
    --nsq=*)
      NSQ="${1#*=}"
      shift
      ;;
    --ws=*)
      WS="${1#*=}"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

CMD="./judge"

if [ -n "$HTTP" ]; then
  CMD="$CMD --http=$HTTP"
fi

if [ -n "$NSQ" ]; then
  CMD="$CMD --nsq=$NSQ"
fi

if [ -n "$WS" ]; then
  CMD="$CMD --ws=$WS"
fi

eval "$CMD & ./go-judge"
