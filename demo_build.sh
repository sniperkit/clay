#!/bin/bash

mkdir -p demo
cd demo
cp ../linux_amd64/clay/clay .
cp ../examples/design.json .
chmod 755 ./clay
(HOST=127.0.0.1 PORT=60000 DB_MODE=file DB_FILE=./clay.db ./clay) &
sleep 5
curl -X PUT '127.0.0.1:60000/designs/present' -H 'Content-Type: application/json' -d @design.json
sleep 2
tar zcvf clay.tgz ./clay ./clay.db

curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.tgz
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.tgz -F content=@clay.tgz -F description="Clay demo"
