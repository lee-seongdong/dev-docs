## Logstash
실시간 파이프라이닝 기능을 가진 데이터수집 엔진  
서로 다른 소스의 데이터(이벤트)를 동적으로 통합하고, 원하는 대상으로 데이터를 정규화 할 수 있음  
플러그인 기반으로 동작하며, 3단계로 처리한다.

### 1. Input
데이터를 수집할 소스를 정의하는 단계. 여러 소스로부터 데이터를 전달받을 수 있다.

#### Plugins
| Plugin          | 설명                         | 사용 예                              |
| --------------- | -------------------------- | --------------------------------- |
| `beats`         | Filebeat 등 Beats로부터 이벤트 수신 | Filebeat → Logstash 파이프라인 연결 시 사용 |
| `file`          | 로컬 파일 읽기                   | 파일 직접 수집 테스트                      |
| `stdin`         | 표준 입력                      | 간단한 테스트용                          |
| `tcp`           | TCP 포트에서 메시지 수신            | 애플리케이션 로그를 직접 TCP로 전송하는 경우        |
| `udp`           | UDP 포트에서 메시지 수신            | syslog, NetFlow 등 수집 시            |
| `http`          | HTTP 요청으로 데이터 수신           | Webhook, 외부 시스템 연동                |
| `kafka`         | Kafka에서 메시지 읽기             | 대규모 분산 로그 처리                      |
| `s3`            | AWS S3에서 파일 읽기             | 배치성 로그 수집                         |
| `syslog`        | Syslog 메시지를 수신             | RFC3164 기반 로그                     |
| `http_poller`   | HTTP로 주기적으로 API 호출         | 외부 REST API 연동 시 사용               |


#### 예시
```conf
input {
  beats {
    port => 5044
  }

  file {
    path => "/var/log/myapp.log"
    start_position => "beginning"
    sincedb_path => "/dev/null"  # 처음부터 읽기
  }

  kafka {
    bootstrap_servers => "localhost:9092"
    topics => ["my-logs"]
    group_id => "logstash-group"
  }
}
```

### 2. Filter
수집한 데이터 가공 및 정규화하는 단계

#### Plugins
| Plugin    | 역할 요약                | 주요 사용 예시                        |
| --------- | ---------------------- | ------------------------------- |
| `if/else` | 조건부 처리                 | 로그 레벨별로 필터나 태그를 다르게 적용          |
| `drop`    | 이벤트 제거                 | 로그 무시   |
| `kv`      | Key-Value 문자열을 필드로 분리  | `key=value` 형태의 로그 파싱         |
| `json`    | JSON 문자열을 필드로 파싱       | `{"key": "value"}` 형식의 메시지를 구조화 |
| `grok`    | 정규표현식 기반 패턴 매칭         | 로그에서 날짜, 레벨, 메시지 등을 추출          |
| `mutate`  | 필드 조작 (변환, 삭제, 소문자화 등) | 타입 변환, 필드 제거, 이름 변경 등           |
| `clone`   | 이벤트 복제                 | 동일한 로그를 다르게 처리하고 싶을 때 사용        |
| `date`    | 날짜 문자열을 타임스탬프로 변환      | 커스텀 타임스탬프 필드를 `@timestamp`로 설정  |


#### 예시
```conf
filter {
  if [status] == "404" {
    drop { }
  }

  if "ERROR" in [level] {
    mutate { add_tag => ["error_log"] }
  }

  kv {
    source => "message"
    field_split => " | "
    value_split => "="
    trim_key => "\""
    trim_value => "\""
  }

  grok {
    match => { "message" => "(?<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[%{DATA:thread}\] %{LOGLEVEL:level} %{JAVACLASS:logger} - %{GREEDYDATA:msg}" }
  }

  mutate {
    convert => { "request_time" => "float" }
    remove_field => [ "host", "path" ]
  }
}
```

> 데이터 가공에 kv나 grok 플러그인을 주로 사용함.  
> grok 패턴 예약어 종류 : https://github.com/elastic/elasticsearch/blob/main/libs/grok/src/main/resources/patterns/legacy/grok-patterns


### 3. Output
가공 및 정규화된 데이터를 최종 목적지로 전달하는 단계. 여러 목적지로 전달할 수 있다.

#### Plugins
| Plugin                | 설명                                  | 주요 사용 예시                |
| --------------------- | ------------------------------------ | ------------------------- |
| `elasticsearch`       | 이벤트를 Elasticsearch 또는 OpenSearch에 전송 | 수집된 로그를 Elasticsearch로 전송 |
| `stdout`              | 이벤트를 표준 출력으로 출력 (디버깅용)               | 파이프라인 테스트/개발 시 확인용        |
| `file`                | 이벤트를 로컬 파일에 저장                       | 로그 백업, 특정 조건 로그 저장        |
| `http`                | HTTP API로 이벤트 전송                     | 외부 API 서버로 전달             |
| `kafka`               | Kafka로 이벤트 전송                        | Kafka 기반 로그 처리 파이프라인 구축   |
| `s3`                  | AWS S3 버킷에 로그 업로드                    | 장기 보관 로그 백업               |
| `loggly`, `datadog` 등 | 서드파티 SaaS 로그 플랫폼에 로그 전송              | SaaS 기반 로그 분석             |

#### 예시
```conf
output {
  if [log_type] == "access-log" {
    opensearch {
      hosts => [ "https://test.com:10200" ]
      index => "access-%{+YYYY.MM.dd}"
      user => "user"
      password => "admin"
      ecs_compatibility => disabled
      cacert => "/etc/logstash/config/certs/root_ca.pem"
    }
  }

  if "error" in [tags] {
    elasticsearch {
      hosts => ["http://localhost:9200"]
      index => "error-logs-%{+YYYY.MM.dd}"
    }
  } else {
    stdout { codec => rubydebug }
  }
}

```