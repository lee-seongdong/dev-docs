# Document
Opensearch에 저장되는 기본 데이터 단위

## Document 구조
```json
{
  "_index": "user-logs",           // 소속 인덱스
  "_id": "unique-id-123",          // 문서의 고유 식별자
  "_source": {                     // 원본 JSON 데이터
    "user": "john",
    "action": "login",
    "timestamp": "2024-01-15T10:30:00Z",
    "ip": "192.168.1.100"
  },
  "_score": 1.5                    // 검색 시 관련도 점수
}
```