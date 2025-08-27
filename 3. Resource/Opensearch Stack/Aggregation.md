# Aggregation
데이터 집계 및 분석을 위한 기능  
검색 결과를 그룹화하고 통계를 계산하여 인사이트를 도출하는 과정

## 1. Metric Aggregations (지표 집계)
```bash
# 기본 통계 (count, sum, avg, min, max)
GET /sales/_search
{
  "size": 0,                    # 검색 결과는 불필요하므로 0
  "aggs": {
    "total_sales": { "sum": { "field": "amount" } },
    "avg_sales": { "avg": { "field": "amount" } },
    "sales_count": { "value_count": { "field": "amount" } }
  }
}

# 고급 통계 (한번에 여러 통계)
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "sales_stats": { "stats": { "field": "amount" } },        # sum, avg, min, max, count
    "sales_extended": { "extended_stats": { "field": "amount" } } # + std_deviation
  }
}

# 백분위수 및 고유값 개수
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "price_percentiles": { "percentiles": { "field": "amount", "percents": [25, 50, 75, 95] } },
    "unique_customers": { "cardinality": { "field": "customer_id" } }
  }
}
```

## 2. Bucket Aggregations (버킷 집계)
```bash
# Terms Aggregation - 그룹별 집계
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "sales_by_category": {
      "terms": { 
        "field": "category",
        "size": 10                # 상위 10개 카테고리
      }
    }
  }
}

# Date Histogram - 시간대별 집계
GET /logs/_search
{
  "size": 0,
  "aggs": {
    "logs_over_time": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1h"    # 1시간 단위
      }
    }
  }
}

# Range Aggregation - 범위별 집계
GET /products/_search
{
  "size": 0,
  "aggs": {
    "price_ranges": {
      "range": {
        "field": "price",
        "ranges": [
          { "to": 100 },
          { "from": 100, "to": 500 },
          { "from": 500 }
        ]
      }
    }
  }
}
```

## 3. 중첩 집계 (Sub-aggregations)
```bash
# 카테고리별 평균 가격 및 판매량
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "categories": {
      "terms": { "field": "category" },
      "aggs": {                          # 중첩 집계
        "avg_price": { "avg": { "field": "price" } },
        "total_sales": { "sum": { "field": "amount" } }
      }
    }
  }
}

# 시간대별 로그 레벨 분포
GET /logs/_search
{
  "size": 0,
  "aggs": {
    "logs_by_hour": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1h"
      },
      "aggs": {
        "log_levels": {
          "terms": { "field": "level" }
        }
      }
    }
  }
}
```

## 4. 필터와 집계 조합
```bash
# 특정 조건의 데이터만 집계
GET /sales/_search
{
  "query": {
    "range": {
      "timestamp": {
        "gte": "2024-01-01",
        "lte": "2024-01-31"
      }
    }
  },
  "size": 0,
  "aggs": {
    "monthly_sales": {
      "sum": { "field": "amount" }
    },
    "top_products": {
      "terms": { "field": "product_id", "size": 5 }
    }
  }
}

# Filters Aggregation - 조건별 집계
GET /logs/_search
{
  "size": 0,
  "aggs": {
    "log_analysis": {
      "filters": {
        "filters": {
          "errors": { "term": { "level": "ERROR" } },
          "warnings": { "term": { "level": "WARN" } },
          "info": { "term": { "level": "INFO" } }
        }
      }
    }
  }
}
```

## 5. 예시
### 5.1. 로그 분석
```bash
# 시간대별 에러 로그 분석
GET /logs/_search
{
  "query": { "term": { "level": "ERROR" } },
  "size": 0,
  "aggs": {
    "errors_over_time": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1h"
      },
      "aggs": {
        "top_error_messages": {
          "terms": { "field": "message.keyword", "size": 3 }
        }
      }
    }
  }
}
```

### 5.2. 비즈니스 데이터 분석
```bash
# 고객별 구매 패턴 분석
GET /orders/_search
{
  "size": 0,
  "aggs": {
    "customer_analysis": {
      "terms": { "field": "customer_id", "size": 100 },
      "aggs": {
        "total_spent": { "sum": { "field": "amount" } },
        "order_count": { "value_count": { "field": "order_id" } },
        "avg_order_value": { "avg": { "field": "amount" } }
      }
    }
  }
}
```


## 6. 성능 최적화
| 최적화 방법 | 설명 | 효과 |
|-------------|------|------|
| `size: 0` | 검색 결과 비활성화 | 메모리 절약 |
| `doc_count` 제한 | `min_doc_count: 1` | 빈 버킷 제거 |
| 필드 캐싱 | `keyword` 필드 사용 | 집계 속도 향상 |
| 필터 우선 적용 | `query` + `aggs` | 집계 대상 데이터 축소 |

**주의사항**:
- **대용량 데이터**: `terms` 집계는 메모리 사용량이 높음
- **카디널리티**: 고유값이 많은 필드는 `composite` 집계 고려
- **정확도**: `terms` 집계는 근사값이므로 정확한 계산이 필요하면 다른 방법 사용


## 관련 문서
- [Query](./Query.md)
- [Index](./Index.md)
- [Index Mappings](./Index%20Mappings.md)
