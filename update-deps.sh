#!/bin/bash

godep update ./...
godep save ./... |& egrep -v 'godep: rewrite: lstat.*no such file or'

echo "Don't forget to git add vendor"
