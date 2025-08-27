# Index Policy
Index의 생명주기를 관리하는 규칙(Action / Transition)의 집합  
인덱스가 언제, 어떤 상태로 전환할지를 설정  
  
Opensearch는 Policy와 **ISM 플러그인**을 통해 Index의 라이프사이클을 관리

## 1. 구성
- [**State**](./Index%20State.md)
- **Action**: 각 상태에서 수행할 작업. rollover, read_only, delete 등
- **Transition**: 다른 상태로 전환할 조건

## 2. 설정
```json
{
  "policy": {
    "description": "정책에 대한 설명",
    "default_state": "시작 상태",
    "states": [
      {
        "name": "상태명1",
        "actions": [/* 해당 상태에서 수행할 액션들 */],
        "transitions": [/* 다음 상태로 전환 조건들 */]
      },
      {
        "name": "상태명2",
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
    "description": "로그 데이터 생명주기",
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
            "state_name": "delete",
            "conditions": {
              "min_index_age": "30d"
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


## 관련 문서
- [Index](./Index.md)
- [Index State](./Index%20State.md)