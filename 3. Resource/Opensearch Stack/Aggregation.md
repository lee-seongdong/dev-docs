# Aggregation
데이터 집계 및 분석을 위한 강력한 기능

## 집계 기본 개념

### 집계 구조
```json
{
  "aggs": {
    "aggregation_name": {
      "aggregation_type": {
        "field": "field_name"
      }
    }
  }
}
```

### 집계 타입 분류
- **Metric Aggregations**: 수치 계산 (sum, avg, max, min 등)
- **Bucket Aggregations**: 그룹화 (terms, date_histogram, range 등)
- **Pipeline Aggregations**: 다른 집계 결과를 입력으로 사용

## Metric 집계

### 기본 수치 집계
```bash
# 여러 metric 집계 한번에
GET /sales/_search
{
  "size": 0,                           // 문서는 반환하지 않고 집계만
  "aggs": {
    "total_sales": {
      "sum": {
        "field": "amount"
      }
    },
    "avg_sales": {
      "avg": {
        "field": "amount"
      }
    },
    "max_sales": {
      "max": {
        "field": "amount"
      }
    },
    "min_sales": {
      "min": {
        "field": "amount"
      }
    },
    "sales_count": {
      "value_count": {
        "field": "amount"
      }
    }
  }
}
```

### Stats 집계 (한번에 여러 통계)
```bash
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "sales_stats": {
      "stats": {
        "field": "amount"               // count, min, max, avg, sum 한번에
      }
    },
    "extended_stats": {
      "extended_stats": {
        "field": "amount"               // 추가로 variance, std_deviation 등
      }
    }
  }
}
```

### Percentiles 집계
```bash
GET /response_times/_search
{
  "size": 0,
  "aggs": {
    "response_time_percentiles": {
      "percentiles": {
        "field": "response_time",
        "percents": [50, 95, 99]        // 50%, 95%, 99% 백분위수
      }
    },
    "response_time_ranks": {
      "percentile_ranks": {
        "field": "response_time",
        "values": [100, 500, 1000]      // 100ms, 500ms, 1000ms의 백분위 순위
      }
    }
  }
}
```

### Cardinality 집계 (고유값 개수)
```bash
GET /user-logs/_search
{
  "size": 0,
  "aggs": {
    "unique_users": {
      "cardinality": {
        "field": "user.keyword"         // 고유 사용자 수
      }
    },
    "unique_ips": {
      "cardinality": {
        "field": "ip"
      }
    }
  }
}
```

## Bucket 집계

### 1. **Terms 집계** (카테고리별 그룹화)
```bash
GET /user-logs/_search
{
  "size": 0,
  "aggs": {
    "top_users": {
      "terms": {
        "field": "user.keyword",
        "size": 10,                     // 상위 10개만
        "order": { "_count": "desc" }   // 문서 수 내림차순
      }
    },
    "actions_by_count": {
      "terms": {
        "field": "action.keyword",
        "min_doc_count": 100           // 최소 100개 문서가 있는 그룹만
      }
    }
  }
}
```

### 2. **Date Histogram** (시간대별 집계)
```bash
GET /user-logs/_search
{
  "size": 0,
  "aggs": {
    "logs_over_time": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1h",      // 1시간 간격
        "format": "yyyy-MM-dd HH:mm",
        "min_doc_count": 0              // 문서가 없는 시간대도 포함
      }
    },
    "daily_logs": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1d",
        "time_zone": "Asia/Seoul"       // 시간대 설정
      }
    }
  }
}
```

### 3. **Range 집계** (숫자 범위별 그룹화)
```bash
GET /products/_search
{
  "size": 0,
  "aggs": {
    "price_ranges": {
      "range": {
        "field": "price",
        "ranges": [
          { "to": 100 },                // 100 미만
          { "from": 100, "to": 500 },   // 100-500
          { "from": 500 }               // 500 이상
        ]
      }
    },
    "age_groups": {
      "range": {
        "field": "age",
        "ranges": [
          { "key": "young", "to": 30 },
          { "key": "middle", "from": 30, "to": 60 },
          { "key": "senior", "from": 60 }
        ]
      }
    }
  }
}
```

### 4. **Date Range 집계**
```bash
GET /events/_search
{
  "size": 0,
  "aggs": {
    "date_ranges": {
      "date_range": {
        "field": "date",
        "format": "yyyy-MM-dd",
        "ranges": [
          { "key": "last_week", "from": "now-7d", "to": "now" },
          { "key": "last_month", "from": "now-1M", "to": "now-7d" }
        ]
      }
    }
  }
}
```

### 5. **Histogram 집계** (숫자 간격별)
```bash
GET /response_times/_search
{
  "size": 0,
  "aggs": {
    "response_time_distribution": {
      "histogram": {
        "field": "response_time",
        "interval": 50                  // 50ms 간격으로 분포
      }
    }
  }
}
```

## 중첩 집계 (Sub-Aggregations)

### Bucket 내부의 Metric 집계
```bash
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "sales_by_category": {
      "terms": {
        "field": "category.keyword"
      },
      "aggs": {                        // 카테고리별 하위 집계
        "total_revenue": {
          "sum": {
            "field": "amount"
          }
        },
        "avg_price": {
          "avg": {
            "field": "amount"
          }
        },
        "sales_stats": {
          "stats": {
            "field": "amount"
          }
        }
      }
    }
  }
}
```

### 다중 레벨 중첩
```bash
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "by_category": {
      "terms": {
        "field": "category.keyword"
      },
      "aggs": {
        "by_month": {
          "date_histogram": {
            "field": "date",
            "calendar_interval": "1M"
          },
          "aggs": {
            "monthly_revenue": {
              "sum": {
                "field": "amount"
              }
            },
            "by_salesperson": {
              "terms": {
                "field": "salesperson.keyword",
                "size": 5
              },
              "aggs": {
                "person_revenue": {
                  "sum": {
                    "field": "amount"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

## Pipeline 집계

### 1. **Bucket Selector** (조건부 버킷 필터링)
```bash
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "monthly_sales": {
      "date_histogram": {
        "field": "date",
        "calendar_interval": "1M"
      },
      "aggs": {
        "total_sales": {
          "sum": {
            "field": "amount"
          }
        },
        "high_sales_months": {           // 매출이 10000 이상인 월만 필터링
          "bucket_selector": {
            "buckets_path": {
              "sales": "total_sales"
            },
            "script": "params.sales > 10000"
          }
        }
      }
    }
  }
}
```

### 2. **Moving Average** (이동 평균)
```bash
GET /metrics/_search
{
  "size": 0,
  "aggs": {
    "daily_values": {
      "date_histogram": {
        "field": "date",
        "calendar_interval": "1d"
      },
      "aggs": {
        "daily_sum": {
          "sum": {
            "field": "value"
          }
        },
        "moving_avg": {
          "moving_avg": {
            "buckets_path": "daily_sum",
            "window": 7,                 // 7일 이동 평균
            "model": "simple"
          }
        }
      }
    }
  }
}
```

### 3. **Derivative** (변화율)
```bash
GET /metrics/_search
{
  "size": 0,
  "aggs": {
    "daily_sales": {
      "date_histogram": {
        "field": "date",
        "calendar_interval": "1d"
      },
      "aggs": {
        "total_sales": {
          "sum": {
            "field": "amount"
          }
        },
        "sales_change": {               // 전일 대비 변화량
          "derivative": {
            "buckets_path": "total_sales"
          }
        }
      }
    }
  }
}
```

### 4. **Cumulative Sum** (누적 합계)
```bash
GET /sales/_search
{
  "size": 0,
  "aggs": {
    "monthly_sales": {
      "date_histogram": {
        "field": "date",
        "calendar_interval": "1M"
      },
      "aggs": {
        "monthly_total": {
          "sum": {
            "field": "amount"
          }
        },
        "cumulative_sales": {           // 누적 매출
          "cumulative_sum": {
            "buckets_path": "monthly_total"
          }
        }
      }
    }
  }
}
```

## 복합 집계 쿼리

### 대시보드용 종합 분석
```bash
GET /user-logs/_search
{
  "size": 0,
  "query": {
    "range": {
      "timestamp": {
        "gte": "now-7d"
      }
    }
  },
  "aggs": {
    "overview": {
      "global": {},                     // 쿼리 필터 무시하고 전체 데이터
      "aggs": {
        "total_logs": {
          "value_count": {
            "field": "timestamp"
          }
        }
      }
    },
    "recent_activity": {
      "date_histogram": {
        "field": "timestamp",
        "calendar_interval": "1h"
      },
      "aggs": {
        "unique_users": {
          "cardinality": {
            "field": "user.keyword"
          }
        },
        "error_rate": {
          "filter": {
            "term": {
              "level.keyword": "ERROR"
            }
          },
          "aggs": {
            "error_count": {
              "value_count": {
                "field": "level"
              }
            }
          }
        }
      }
    },
    "top_errors": {
      "filter": {
        "term": {
          "level.keyword": "ERROR"
        }
      },
      "aggs": {
        "error_messages": {
          "terms": {
            "field": "message.keyword",
            "size": 10
          }
        }
      }
    },
    "user_behavior": {
      "terms": {
        "field": "action.keyword"
      },
      "aggs": {
        "unique_users": {
          "cardinality": {
            "field": "user.keyword"
          }
        },
        "avg_response_time": {
          "avg": {
            "field": "response_time"
          }
        }
      }
    }
  }
}
```

## 성능 최적화

### 1. **필드 최적화**
```bash
# doc_values 활용 (기본적으로 활성화)
# keyword, numeric, date 필드는 자동으로 doc_values 사용
# text 필드는 fielddata 필요 (메모리 사용량 증가)

PUT /my-index
{
  "mappings": {
    "properties": {
      "category": {
        "type": "keyword"               // doc_values 자동 활성화
      },
      "description": {
        "type": "text",
        "fielddata": true               // 집계 시 필요하지만 메모리 사용량 증가
      }
    }
  }
}
```

### 2. **집계 캐시 활용**
```bash
# 동일한 집계 결과 캐시
GET /logs/_search?request_cache=true
{
  "size": 0,
  "aggs": {
    "cached_agg": {
      "terms": {
        "field": "status.keyword"
      }
    }
  }
}
```

### 3. **샤드 크기 고려**
```bash
# 너무 많은 샤드는 집계 성능 저하
# 집계 대상 데이터가 많은 경우 샤드 수 최적화 필요

# 집계 성능 모니터링
GET /_cat/nodes?v&h=name,search.query_total,search.query_time_in_millis
```

### 4. **정확도 vs 성능 트레이드오프**
```bash
GET /large-index/_search
{
  "size": 0,
  "aggs": {
    "approximate_cardinality": {
      "cardinality": {
        "field": "user_id",
        "precision_threshold": 1000    // 정확도 조정 (기본: 3000)
      }
    },
    "sampled_terms": {
      "sampler": {
        "shard_size": 1000            // 샘플링으로 성능 향상
      },
      "aggs": {
        "top_categories": {
          "terms": {
            "field": "category.keyword"
          }
        }
      }
    }
  }
}
```

## 집계 결과 분석

### 결과 구조
```json
{
  "aggregations": {
    "sales_by_category": {
      "doc_count_error_upper_bound": 0,
      "sum_other_doc_count": 0,
      "buckets": [
        {
          "key": "Electronics",
          "doc_count": 1234,
          "total_revenue": {
            "value": 567890.5
          }
        }
      ]
    }
  }
}
```

### 주요 메트릭
- **doc_count**: 해당 버킷의 문서 수
- **doc_count_error_upper_bound**: 문서 수 오차 상한선
- **sum_other_doc_count**: 결과에 포함되지 않은 나머지 문서 수

### 집계 디버깅
```bash
# Profile API로 집계 성능 분석
GET /logs/_search
{
  "profile": true,
  "size": 0,
  "aggs": {
    "debug_agg": {
      "terms": {
        "field": "category.keyword"
      }
    }
  }
}
```
