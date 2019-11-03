go run genenum.go -typename=AchieveType -packagename=achievetype -basedir=. -statstype=int
goimports -w achievetype/achievetype_gen.go
goimports -w achievetype_stats/achievetype_stats_gen.go
