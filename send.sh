#!/bin/sh

curl -v `hostname --fqdn`:8585/person --data '{"id":42}'
