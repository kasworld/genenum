#!/usr/bin/env bash

# go run genenum.go -typename=AchieveType -packagename=achievetype -basedir=. -statstype=int
go run genenum.go -typename=AchieveType -packagename=achievetype -basedir=. -flagtype=uint32 -statstype=int -verbose

goimports -w achievetype/achievetype_gen.go
goimports -w achievetype/achievetype_flag_gen.go
goimports -w achievetype_stats/achievetype_stats_gen.go
