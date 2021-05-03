#!/usr/bin/env bash

set -eo pipefail

hugo -e production && aws s3 sync public/ s3://inge4pres-website
aws cloudfront create-invalidation --distribution-id E2THBBYGNS57TT --paths "/*"
