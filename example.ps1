
echo "make enum, flag, vector package"
go run genenum.go -typename=AchieveType -packagename=achievetype -basedir="." -flagtype=uint32 -vectortype="int,float64" -verbose

goimports -w achievetype
goimports -w achievetype_flag
goimports -w achievetype_vector_float64
goimports -w achievetype_vector_int


echo "***** too many enum for flag package example *****"
go run genenum.go -typename=BigEnumType -packagename=bigenumtype -basedir="." -flagtype=uint64 -vectortype=int -verbose

goimports -w bigenumtype
goimports -w bigenumtype_flag
goimports -w bigenumtype_vector

