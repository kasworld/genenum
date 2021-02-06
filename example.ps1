
echo "make enum, flag, vector package"
go run genenum.go -typename=AchieveType -packagename=achievetype -basedir="." -flagtype=uint32 -vectortype=int -verbose

goimports -w achievetype
goimports -w achievetype_flag
goimports -w achievetype_vector


echo "***** too many enum for flag package test *****"
go run genenum.go -typename=BigEnumType -packagename=bigenumtype -basedir="." -flagtype=uint64 -vectortype=int -verbose

goimports -w bigenumtype
goimports -w bigenumtype_flag
goimports -w bigenumtype_vector

