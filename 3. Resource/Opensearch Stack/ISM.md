# ISM (Index State Management)
인덱스의 생명주기를 자동으로 관리하는 플러그인

## 1. State
인덱스가 생명주기 동안 거치는 **논리적 상태**를 의미  
각 상태는 데이터의 **접근 빈도, 성능 요구사항, 저장 비용**에 따라 구분  
주로 데이터의 온도에 따라 Hot, Warm, Cold, Frozen 등으로 분류함  
(kafka, hadoop, cloud storage 등에서 이미 사용중인 명칭이고, 직관적이기 때문)

### 주요 Index State 종류
#### 1. **Hot**
- **용도**: 최신 데이터 저장 및 활발한 색인/검색
- **하드웨어**: SSD, 고성능 CPU, 충분한 메모리
- **샤드**: 더 많은 샤드로 분산 처리
- **설정**: `node.attr.tier: hot`

#### 2. **Warm**
- **용도**: 이전 데이터, 가끔 검색하는 데이터
- **하드웨어**: SATA 디스크, 적당한 성능
- **샤드**: 적은 수의 샤드로 통합
- **설정**: `node.attr.tier: warm`

#### 3. **Cold**
- **용도**: 아카이브 데이터, 거의 검색하지 않음
- **하드웨어**: 저비용 스토리지
- **샤드**: 최소한의 샤드
- **설정**: `node.attr.tier: cold`

### State별 최적화 전략
#### **Hot State 최적화**
- **샤드 수**: 더 많은 샤드로 분산 처리
- **복제본**: 읽기 성능을 위한 충분한 복제본
- **새로고침**: 빠른 검색을 위한 1초 주기
- **하드웨어**: NVMe SSD, 고성능 CPU

#### **Warm State 최적화**  
- **샤드 통합**: Force Merge로 세그먼트 최적화
- **압축**: 스토리지 효율성 증대
- **복제본 감소**: 비용 절약
- **하드웨어**: SATA SSD, 일반 CPU

#### **Cold State 최적화**
- **최대 압축**: best_compression 코덱 사용
- **최소 복제본**: 0개 또는 1개
- **인덱스 통합**: Shrink API로 샤드 수 축소
- **하드웨어**: HDD, 저사양 CPU

## 2. Policy
ISM에서 인덱스 생명주기를 관리하는 규칙과 액션(**인덱스가 언제, 어떤 상태로 전환할지**)을 정의하는 설정  
Policy는 **State(상태)**, **Action(액션)**, **Transition(전환 조건)**으로 구성

### 최적화
#### 1. **Rollover 크기 최적화**
- **작은 인덱스**: 관리 오버헤드 증가
- **큰 인덱스**: 성능 저하, 복구 시간 증가
- **권장 크기**: 10GB-50GB

#### 2. **Transition 시점 최적화**
- **너무 빠른 전환**: 리소스 낭비
- **너무 늦은 전환**: Hot tier 부하 증가
- **데이터 접근 패턴 분석** 필요


### 설정
json 파일로 설정하며, Opensearch Dashboard 에서도 설정가능
구조: 
```json
{
  "policy": {
    "description": "정책에 대한 설명",
    "default_state": "시작 상태",
    "states": [
      {
        "name": "상태명",
        "actions": [/* 해당 상태에서 수행할 액션들 */],
        "transitions": [/* 다음 상태로 전환 조건들 */]
      }
    ]
  }
}
```

예시: 
```json
{
  "policy": {
    "description": "로그 데이터 생명주기 관리",
    "default_state": "hot",
    "states": [
      {
        "name": "hot",
        "actions": [
          {
            "rollover": {
              "min_index_age": "1d",
              "min_size": "5gb"
            }
          }
        ],
        "transitions": [
          {
            "state_name": "warm",
            "conditions": {
              "min_index_age": "7d"
            }
          }
        ]
      },
      {
        "name": "warm", 
        "actions": [
          {
            "allocation": {
              "require": { "tier": "warm" },
              "include": {},
              "exclude": {}
            }
          },
          {
            "force_merge": {
              "max_num_segments": 1
            }
          }
        ],
        "transitions": [
          {
            "state_name": "cold",
            "conditions": {
              "min_index_age": "30d"
            }
          }
        ]
      },
      {
        "name": "cold",
        "actions": [
          {
            "allocation": {
              "require": { "tier": "cold" }
            }
          },
          {
            "force_merge": {
              "max_num_segments": 1
            }
          }
        ],
        "transitions": [
          {
            "state_name": "delete",
            "conditions": {
              "min_index_age": "90d"
            }
          }
        ]
      },
      {
        "name": "delete",
        "actions": [
          {
            "delete": {}
          }
        ]
      }
    ]
  }
}
```



## 3. 활용
### 주요 적용 시나리오
1. **로그 데이터 관리**: 시간이 지날수록 접근 빈도가 감소
2. **시계열 데이터**: 메트릭, 모니터링 데이터
3. **규제 요구사항**: 일정 기간 보관 후 삭제
4. **비용 최적화**: 스토리지 비용 절약
5. **성능 최적화**: Hot tier의 부하 분산

### ISM을 위한 노드 설정
#### Tier별 노드 구성
```yaml
# Hot 노드 설정
cluster.name: opensearch-cluster
node.name: opensearch-hot-1
node.roles: [cluster_manager, data, ingest]
node.attr.tier: hot

# Warm 노드 설정  
cluster.name: opensearch-cluster
node.name: opensearch-warm-1
node.roles: [data]
node.attr.tier: warm

# Cold 노드 설정
cluster.name: opensearch-cluster  
node.name: opensearch-cold-1
node.roles: [data]
node.attr.tier: cold
```

#### 클러스터 할당 설정
```yaml
# Hot-Warm-Cold 할당 설정
cluster.routing.allocation.awareness.attributes: tier
cluster.routing.allocation.awareness.force.tier.values: hot, warm, cold
```