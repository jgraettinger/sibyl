#!/bin/bash
go build main && ./main | dot -T svg > /tmp/foo.svg && /opt/google/chrome/chrome /tmp/foo.svg
