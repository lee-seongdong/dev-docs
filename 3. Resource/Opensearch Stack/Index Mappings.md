# Index Mappings
인덱스에 저장될 문서의 필드 타입과 구조를 정의하는 스키마  
검색 성능과 정확도에 직접적인 영향을 준다.  
인덱스 생성 시 개별적으로 설정할 수도 있고, [Index Template](./Index%20Template.md)의 설정을 따를 수도 있다.


## 1. 매핑 유형
- **정적 매핑**: 인덱스 생성 시 필드 타입을 명시적으로 정의 (권장)
- **동적 매핑**: 문서 삽입 시 필드 타입을 자동으로 추론
  - 장점: 빠른 개발, 유연성
  - 단점: 예상치 못한 타입 추론, 성능 이슈 가능

## 2. 필드 타입
### 2.1. 텍스트 타입
```json
{
  "properties": {
    "title": { "type": "text", "analyzer": "standard" }, # 전문검색. 분석기 적용
    "status": { "type": "keyword" } # 분석기 미적용. 정렬/집계 가능
  }
}
```

### 2.2. 숫자 타입
```json
{
  "properties": {
    "age": { "type": "integer" },
    "price": { "type": "float" },
    "timestamp": { "type": "long" }
  }
}
```

### 2.3. 날짜 타입
```json
{
  "properties": {
    "created_at": { 
      "type": "date",
      "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
    }
  }
}
```

### 2.4. 불린 타입
```json
{
  "properties": {
    "is_active": { "type": "boolean" }
  }
}
```

### 2.5. 객체 타입
```json
{
  "properties": {
    "user": {
      "type": "object",
      "properties": {
        "name": { "type": "text" },
        "email": { "type": "keyword" }
      }
    },
    "tags": {
      "type": "nested",
      "properties": {
        "name": { "type": "keyword" },
        "score": { "type": "float" }
      }
    }
  }
}
```

### 2.6. 특수 타입
```json
{
  "properties": {
    "location": { "type": "geo_point" },
    "ip_address": { "type": "ip" },
    "binary_data": { "type": "binary" }
  }
}
```

## 3. 분석기
### 3.1. 표준 분석기
```json
{
  "properties": {
    "content": {
      "type": "text",
      "analyzer": "standard"        // 기본 분석기
    }
  }
}
```

### 3.2. 언어별 분석기
```json
{
  "properties": {
    "title_en": { "type": "text", "analyzer": "english" },
    "title_ko": { "type": "text", "analyzer": "cjk" }
  }
}
```

### 3.3. 커스텀 분석기
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_analyzer": {
          "tokenizer": "standard",
          "filter": ["lowercase", "stop"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "content": { "type": "text", "analyzer": "my_analyzer" }
    }
  }
}
```


## 4. 관리
### 4.1. 매핑 생성 및 조회
```bash
# 인덱스 생성(정적 매핑)
PUT /my-index
{
  "mappings": {
    "properties": {
      "title": { "type": "text" },
      "price": { "type": "float" }
    }
  }
}

# 매핑 조회
GET /my-index/_mapping

# 특정 필드 매핑 조회
GET /my-index/_mapping/field/title
```

### 4.2. 매핑 업데이트
```bash
# 새 필드 추가 (기존 필드 변경 불가)
PUT /my-index/_mapping
{
  "properties": {
    "new_field": { "type": "keyword" }
  }
}
```

### 4.3. 동적 매핑 제어
```json
{
  "mappings": {
    "dynamic": "strict",          // 새 필드 추가 차단
    "properties": {
      "title": { "type": "text" }
    }
  }
}
```

**동적 매핑 옵션**:
- `true`: 새 필드 자동 추가 (기본값)
- `false`: 새 필드 무시
- `strict`: 새 필드 추가 시 오류 발생


## 5. 최적화
### 5.1. 성능 최적화 팁
- **불필요한 분석 피하기**: 검색하지 않는 필드는 `keyword` 사용
- **적절한 분석기 선택**: 언어와 용도에 맞는 분석기
- **중첩 객체 최소화**: `nested` 타입은 성능에 영향

### 5.2. 일반적인 매핑 패턴
```json
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },               // ID는 keyword
      "title": { "type": "text" },               // 검색용 텍스트
      "status": { "type": "keyword" },           // 필터링용
      "created_at": { "type": "date" },          // 날짜/시간
      "price": { "type": "scaled_float", "scaling_factor": 100 }, // 가격
      "tags": { "type": "keyword" },             // 다중값 키워드
      "metadata": { "enabled": false }          // 검색 불필요 시 비활성화
    }
  }
}
```

### 5.3. 주의사항
- **매핑 변경 제약**: 기존 필드 타입은 변경 불가 (reindex 필요)
- **동적 매핑 위험성**: 예상치 못한 타입 추론으로 인한 문제
- **분석기 일관성**: 인덱싱과 검색 시 동일한 분석기 사용 권장


## 관련 문서
- [Index](./Index.md)
- [Index Settings](./Index%20Settings.md)
- [Index Template](./Index%20Template.md)
