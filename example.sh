#!/usr/bin/env bash

go run genenum.go -typename=AchieveType -packagename=achievetype -basedir=. -flagtype=uint32 -vectortype=int -verbose

goimports -w achievetype/achievetype_gen.go
goimports -w achievetype_flag/achievetype_flag_gen.go
goimports -w achievetype_vector/achievetype_vector_gen.go
