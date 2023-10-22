#!/bin/bash

cd "$(dirname "$0")/.." || exit

sql-migrate down "$@"