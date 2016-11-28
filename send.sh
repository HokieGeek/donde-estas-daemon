#!/bin/sh

curl `hostname --fqdn`:8585/person --data '{"id":42}'
