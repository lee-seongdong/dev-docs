# Index Mappings
인덱스의 필드 타입과 구조를 정의하는 스키마

## 매핑 기본 개념

### 매핑이란?
- 인덱스에 저장될 문서의 필드 타입과 구조를 정의
- RDBMS의 테이블 스키마와 유사한 개념
- 검색 성능과 정확성에 직접적인 영향

### 동적 vs 정적 매핑
```json
// 동적 매핑: 자동으로 필드 타입 추론
{
  "name": "john",           // → text, keyword
  "age": 25,                // → long  
  "created": "2024-01-15"   // → date
}

// 정적 매핑: 명시적으로 타입 정의
{
  "mappings": {
    "properties": {
      "name": { "type": "keyword" },
      "age": { "type": "integer" },
      "created": { "type": "date", "format": "yyyy-MM-dd" }
    }
  }
}
```

## 주요 필드 타입

### 1. **Text 타입**
전문 검색용 텍스트 필드
```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "standard",           // 분석기 설정
        "search_analyzer": "standard",    // 검색 시 분석기
        "fields": {
          "keyword": {                    // Multi-field: keyword 서브필드
            "type": "keyword"
          }
        }
      }
    }
  }
}
```

### 2. **Keyword 타입**
정확한 값 매치, 집계, 정렬용
```json
{
  "mappings": {
    "properties": {
      "status": {
        "type": "keyword",
        "ignore_above": 256               // 256자 초과 시 무시
      },
      "user_id": {
        "type": "keyword",
        "store": true                     // _source와 별도 저장
      }
    }
  }
}
```

### 3. **Date 타입**
날짜/시간 데이터 처리
```json
{
  "mappings": {
    "properties": {
      "created_at": {
        "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||epoch_millis"
      },
      "birthday": {
        "type": "date",
        "format": "yyyy-MM-dd"
      },
      "timestamp": {
        "type": "date",
        "format": "epoch_second"          // Unix timestamp
      }
    }
  }
}
```

### 4. **Numeric 타입**
숫자 데이터 처리
```json
{
  "mappings": {
    "properties": {
      "age": { "type": "integer" },
      "price": { "type": "float" },
      "score": { "type": "double" },
      "count": { "type": "long" },
      "percentage": { "type": "scaled_float", "scaling_factor": 100 }
    }
  }
}
```

### 5. **Boolean & IP 타입**
```json
{
  "mappings": {
    "properties": {
      "is_active": { "type": "boolean" },
      "client_ip": { "type": "ip" },
      "server_ip": { "type": "ip" }
    }
  }
}
```

### 6. **Object & Nested 타입**
중첩 객체 처리
```json
{
  "mappings": {
    "properties": {
      // Object 타입 (기본)
      "user": {
        "type": "object",
        "properties": {
          "name": { "type": "keyword" },
          "email": { "type": "keyword" }
        }
      },
      
      // Nested 타입 (배열 객체 정확한 검색)
      "comments": {
        "type": "nested",
        "properties": {
          "author": { "type": "keyword" },
          "message": { "type": "text" },
          "date": { "type": "date" }
        }
      }
    }
  }
}
```

## 분석기(Analyzer) 설정

### 기본 분석기
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "stop"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "content": {
        "type": "text",
        "analyzer": "my_analyzer"
      }
    }
  }
}
```

### 한국어 분석기 (nori)
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "korean_analyzer": {
          "type": "custom",
          "tokenizer": "nori_tokenizer",
          "filter": ["nori_part_of_speech", "lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "korean_text": {
        "type": "text",
        "analyzer": "korean_analyzer"
      }
    }
  }
}
```

## Multi-field 매핑

### 텍스트 + 키워드 조합
```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",                   // 전문 검색용
        "fields": {
          "keyword": {                    // 정확한 매치, 집계용
            "type": "keyword",
            "ignore_above": 256
          },
          "length": {                     // 길이 기반 검색
            "type": "token_count",
            "analyzer": "standard"
          }
        }
      }
    }
  }
}
```

### 숫자 + 키워드 조합
```json
{
  "mappings": {
    "properties": {
      "product_code": {
        "type": "keyword",                // 정확한 검색용
        "fields": {
          "numeric": {                    // 숫자 연산용
            "type": "long",
            "ignore_malformed": true
          }
        }
      }
    }
  }
}
```

## 커스텀 매핑 전략

### 로그 데이터 매핑
```json
{
  "mappings": {
    "properties": {
      "@timestamp": { "type": "date" },
      "level": { "type": "keyword" },
      "logger": { "type": "keyword" },
      "thread": { "type": "keyword" },
      "message": {
        "type": "text",
        "fields": {
          "keyword": { "type": "keyword", "ignore_above": 1024 }
        }
      },
      "exception": {
        "type": "object",
        "properties": {
          "class": { "type": "keyword" },
          "message": { "type": "text" },
          "stacktrace": { 
            "type": "text",
            "index": false                // 검색하지 않음, 저장만
          }
        }
      },
      "request": {
        "type": "object", 
        "properties": {
          "method": { "type": "keyword" },
          "url": { "type": "keyword" },
          "status": { "type": "integer" },
          "duration": { "type": "float" },
          "user_agent": { 
            "type": "text",
            "fields": {
              "keyword": { "type": "keyword", "ignore_above": 512 }
            }
          }
        }
      }
    }
  }
}
```

### 전자상거래 데이터 매핑
```json
{
  "mappings": {
    "properties": {
      "product_id": { "type": "keyword" },
      "name": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" },
          "ngram": {
            "type": "text",
            "analyzer": "autocomplete_analyzer"
          }
        }
      },
      "category": {
        "type": "object",
        "properties": {
          "id": { "type": "keyword" },
          "name": { "type": "keyword" },
          "path": { "type": "keyword" }    // "Electronics > Computers > Laptops"
        }
      },
      "price": {
        "type": "object",
        "properties": {
          "amount": { "type": "double" },
          "currency": { "type": "keyword" }
        }
      },
      "attributes": {
        "type": "nested",                 // 동적 속성들
        "properties": {
          "name": { "type": "keyword" },
          "value": { "type": "keyword" }
        }
      },
      "tags": { "type": "keyword" },      // 배열 태그
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}
```

## 매핑 변경 및 Reindex

### 매핑 변경 제약사항
- 기존 필드의 타입 변경 불가
- 새 필드 추가는 가능
- 기존 필드 삭제 불가

### 새 필드 추가
```bash
# 새 필드 매핑 추가
PUT /my-index/_mapping
{
  "properties": {
    "new_field": {
      "type": "keyword"
    }
  }
}
```

### Reindex를 통한 매핑 변경
```bash
# 1. 새 인덱스 생성 (새 매핑)
PUT /new-index
{
  "mappings": {
    "properties": {
      "field1": { "type": "text" },      // 기존: keyword → 변경: text
      "field2": { "type": "date" }       // 기존 필드 유지
    }
  }
}

# 2. 데이터 복사
POST /_reindex
{
  "source": { "index": "old-index" },
  "dest": { "index": "new-index" }
}

# 3. 별칭 교체
POST /_aliases
{
  "actions": [
    { "remove": { "index": "old-index", "alias": "my-index" }},
    { "add": { "index": "new-index", "alias": "my-index" }}
  ]
}
```

## 매핑 템플릿

### 동적 템플릿
```json
{
  "mappings": {
    "dynamic_templates": [
      {
        "strings_as_keywords": {
          "match_mapping_type": "string",
          "mapping": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      {
        "integers_as_longs": {
          "match_mapping_type": "long",
          "mapping": {
            "type": "long"
          }
        }
      }
    ]
  }
}
```

### 인덱스 템플릿
```json
PUT /_index_template/logs-template
{
  "index_patterns": ["logs-*"],
  "template": {
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "level": { "type": "keyword" },
        "message": {
          "type": "text",
          "fields": {
            "keyword": { "type": "keyword", "ignore_above": 1024 }
          }
        }
      }
    }
  },
  "priority": 100
}
```

## 매핑 확인 및 분석

### 매핑 정보 확인
```bash
# 전체 매핑 확인
GET /my-index/_mapping

# 특정 필드 매핑 확인
GET /my-index/_mapping/field/field_name

# 분석기 테스트
GET /my-index/_analyze
{
  "field": "title",
  "text": "검색할 텍스트"
}
```

### 매핑 최적화 팁
1. **불필요한 필드 인덱싱 비활성화**: `"index": false`
2. **큰 텍스트는 저장만**: `"store": true, "index": false`
3. **집계 전용 필드**: `"doc_values": true, "index": false`
4. **동적 매핑 제어**: `"dynamic": "strict"` 또는 `"dynamic": false`
