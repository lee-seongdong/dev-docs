## Index
데이터를 저장하는 논리적 공간  
데이터베이스의 테이블과 유사한 성격을 가지며, 각 인덱스는 여러 샤드로 구성됨

### Index Policy
인덱스의 데이터 보관 기간, rollover, 삭제 등 인덱스를 자동으로 관리하는 정책  
> Elasticsearch에서는 `Index Lifecycle Management(ILM)` 라고 부르고, Opensearch에서는 `Index State Management(ISM)` 으로 불림  

#### Index State
Index Policy에 정의된 단계를 의미
- hot: 데이터가 자주 색인되고 조회되는 상태
- warm: 데이터가 종종 검색되지만, 색인은 되지않는 상태
- cold: 데이터가 거의 사용되지 않는 상태
- delete: 삭제 예정 상태
