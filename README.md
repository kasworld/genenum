# generate enumeration and statistics 

입력 data 기반으로 

enum package를 만든다. 

선택적으로 통계 package 를 만들 수도 있다. 

## 사용법 

실행 인자 

typename : 만들 enum의 type 

basedir : 만들어질 package들일 속할 폴더 

packagename : 만들어질 패키지이름(==폴더이름) 동시에 읽어 들일 data 파일 이름, 
    
    packagename.enum 파일을 basedir에서 읽어 들인다. 

statstype : 통계 패키지를 만들때 사용할 element 타입, 없으면 통계 패키지를 만들지 않는다. 

데이터 파일 형태 

    한줄에 한개의 enum을 정의 
    각 라인의 첫 단어가 enum 이고 space 로 분리된 뒷 부분은 생성된 코드의 comment가 된다. 
    # 으로 시작하는 라인은 무시(comment취급)
    achievetype.enum 참고 

생성이 끝난 코드들은 import code가 제대로 되어 있지 않으니 
goimports 등으로 정리 해주어야 합니다. 

실행하면 goimports 를 해야할 파일 목록을 찍어 줍니다. 

    example.sh를 실행한 결과 
    goimports -w achievetype/achievetype_gen.go
    goimports -w achievetype_stats/achievetype_stats_gen.go
