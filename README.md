# generate enumeration and vector(statistics)

입력 data 기반으로 

enum 패키지를 만든다. 

선택적으로 bit flag 패키지를 만들 수 있다.  

선택적으로 vector(통계) 패키지 를 만들 수도 있다. 

## 사용법 

실행 인자 

typename : 만들 enum의 type 

basedir : 만들어질 패키지들일 속할 폴더 

packagename : 만들어질 패키지이름(==폴더이름) 동시에 읽어 들일 data 파일 이름, 
    
    packagename.enum 파일을 basedir에서 읽어 들인다. 

vectortype : 벡터 패키지를 만들때 사용할 element 타입, 없으면 벡터 패키지를 만들지 않는다. 

flagtype : bit flag 패키지를 만들때 사용할 타입 enum 갯수보다 큰 bit len을 가져야 한다. 

    uint8, uint16, uint32, uint64 중하나를 추천

verbose : goimports 를 해야할 파일 목록을 찍어 줍니다. 

데이터 파일 형태 

    한줄에 한개의 enum을 정의 
    각 라인의 첫 단어가 enum 이고 space 로 분리된 뒷 부분은 생성된 코드의 comment가 된다. 
    # 으로 시작하는 라인은 무시(comment취급)
    achievetype.enum 참고 

생성이 끝난 코드들은 import code가 제대로 되어 있지 않으니 
goimports 등으로 정리 해주어야 합니다. 

    example.sh를 실행한 결과 
    goimports -w achievetype/achievetype_gen.go
    goimports -w achievetype_flag/achievetype_flag_gen.go
    goimports -w achievetype_vector/achievetype_vector_gen.go
