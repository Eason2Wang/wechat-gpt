#!/bin/bash
pm2 stop main
go build main.go
pm2 start main
pm2 list
