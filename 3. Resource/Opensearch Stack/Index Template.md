# Index Template
인덱스 생성 시 인덱스 패턴에 따라 자동으로 적용되는 설정 및 매핑 템플릿

## 1. 템플릿 예시
```bash
# 템플릿 1
PUT /_index_template/logs-template
{
  "index_patterns": ["logs-*"],
  "template": {
    "settings": {
      "number_of_shards": 3,
      "number_of_replicas": 1,
      "refresh_interval": "5s",
      "codec": "best_compression"
    },
    "mappings": {
      "properties": {
        "timestamp": { "type": "date" },
        "message": { "type": "text" },
        "level": { "type": "keyword" }
      }
    }
  },
  "priority": 100
}

# 템플릿2
PUT /_index_template/dev-template
{
  "index_patterns": ["dev-*"],
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "refresh_interval": "1s"
    },
    "mappings": {
      "properties": {
        "timestamp": { "type": "date" },
        "message": { "type": "text" },
        "level": { "type": "keyword" }
      }
    }
  },
  "priority": 200
}

# 템플릿3
PUT /_index_template/prod-template  
{
  "index_patterns": ["prod-*"],
  "template": {
    "settings": {
      "number_of_shards": 5,
      "number_of_replicas": 2,
      "refresh_interval": "30s",
      "codec": "best_compression",
      "index.routing.allocation.require.tier": "hot"
    },
    "mappings": {
      "properties": {
        "timestamp": { "type": "date" },
        "message": { "type": "text" },
        "level": { "type": "keyword" }
      }
    }
  },
  "priority": 300
}
```

## 2. 템플릿 관리
```bash
# 템플릿 조회
GET /_index_template/logs-template

# 모든 템플릿 조회
GET /_index_template

# 템플릿 삭제
DELETE /_index_template/logs-template

# 템플릿 적용 시뮬레이션
POST /_index_template/_simulate_index/logs-2024.01
```

## 3. 우선순위 및 상속
- **우선순위**: 높은 값이 우선 적용 (기본값: 0)
- **패턴 매칭**: 더 구체적인 패턴이 우선
- **설정 병합**: 여러 템플릿이 매칭되면 설정이 병합됨

```bash
# 우선순위 예시
PUT /_index_template/general-logs
{
  "index_patterns": ["logs-*"],
  "priority": 100,
  "template": { "settings": { "number_of_shards": 1 } }
}

PUT /_index_template/application-logs  
{
  "index_patterns": ["logs-app-*"],
  "priority": 200,  # 더 높은 우선순위
  "template": { "settings": { "number_of_shards": 3 } }
}
```

## 관련 문서
- [Index](./Index.md)