#!/bin/bash -e

# For dev environment only!!!
# requires a mysql client to be installed locally
# creating a decent docker dev environment is out of scope here
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
mysql -h127.0.0.1 -P3307 -u root -ptest bookmarks < $DIR/schema.sql
