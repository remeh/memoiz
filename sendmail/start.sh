#!/bin/bash

export SMTP_HOST=smtp.mailgun.org
export SMTP_LOGIN=postmaster@sandboxe6a94a1aea874ac5ad7e84d2f40d9a23.mailgun.org
export SMTP_PWD=8aef08030a964902a22f0674915a2601
export SMTP_PORT=587
export CONN='host=localhost sslmode=disable user=memoiz dbname=memoiz password=memoiz'
export RES=/Users/remy/docs/code/gopath/src/remy.io/memoiz/resources

go build

rc=$?
if [[ $rc == 0 ]]; then
        ./sendmail
fi;


