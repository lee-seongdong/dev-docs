## Filebeat
파일의 변화를 실시간으로 감지하고, Logstash 또는 Opensearch로 전송하는 경량 agent


### 1. 구성요소
#### - Input
- 설정된 경로를 주기적으로 폴링하여 새 파일이 생겼는지, 기존 파일의 변경이 있는지 확인하는 컴포넌트
- 새 파일이나 변경된 파일을 감지하면, Harvester를 생성함

#### - Harvester
- 파일의 새 라인을 실시간으로 수집하는 컴포넌트
- 각 파일마다 새로운 Harvester가 생성되어 실행되며, 수집된 로그는 이벤트로 만들어 Spooler로 전달
- 파일의 상태는 registry에 저장되어, 로그 유실없이 수집할 수 있음
- `clean_removed`나 `clean_inactivate` 설정으로 registry가 너무 커지지 않게 설정할 수 있음

#### - Queue
- Harvester로 부터 받은 이벤트를 배치처리 하기 위한 컴포넌트
- 설정된 이벤트 수 혹은 일정 시간마다 이벤트들을 Output으로 전달
- 기본적으로 메모리 기반으로 동작하지만, `queue.disk` 설정을 통해 디스크기반으로 동작하게 설정할 수 있음
- 큐가 가득 차는 경우, `백프레셔(Back Pressure)` 메커니즘으로 데이터 손실 방지(파일 커서 중지)
- 설정에 따른 두가지 모드
  - 비동기 모드(default) : 효율적인 배치처리
  - 동기 모드 (`flush.min_events: 1` or `flush.timeout: 0s` 로 설정) : 실시간 처리

#### - Output (Publisher)
- 넘겨받은 이벤트들을 Logstash, Opensearch, Kafka 등으로 전송하는 컴포넌트


### 설정 예시:
```yaml
# Input 설정
filebeat.inputs:
  - type: filestream
    id: my-app-logs    # filestream에서는 고유 ID 필수 (중복 방지)
    paths:
      - /var/log/app/*.log
    ignore_older: 24h
    fields:
      log_type: access-log  # 커스텀 필드 추가
    fields_under_root: true
    
    parsers:
      - multiline:
          type: pattern
          pattern: '^\d{4}-\d{2}-\d{2}'  # 패턴과 일치하는 줄을 새 로그 시작점으로 설정
          negate: true
          match: after
    
    # 파일 관리 설정
    close.on_state_change.inactive: 1h  # 비활성 파일 닫기
    prospector.scanner.check_interval: 10s  # 파일 스캔 주기
    exclude_lines: ['^DEBUG'] # 특정 라인 제외


# Harvester Registry 설정
filebeat.registry:
  clean_removed: true       # 삭제된 파일 registry 정리
  clean_inactive: 24h       # 비활성 파일 registry 정리


# Queue 설정 (메모리 기반 - 기본값)
queue.mem:
  events: 3200
  flush.min_events: 1600
  flush.timeout: 10s


# Output 설정
output.logstash:
  hosts: ["localhost:5044"]


# Filebeat 로깅 설정
logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  keepfiles: 7
```