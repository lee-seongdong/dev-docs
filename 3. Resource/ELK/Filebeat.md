## Filebeat
파일의 변화를 실시간으로 감지하고, Logstash 또는 Opensearch로 전송하는 경량 agent


### 설정 예시:
```yaml
filebeat.inputs:
  - type: log # log가 가장 일반적인 타입. stdin, syslog, container, filestream 등이 있음
    paths:
      - /var/log/app/*.log
    
    multiline.pattern: '^\d{4}-\d{2}-\d{2}'  # 패턴과 일치하는 줄을 새 로그 시작점으로 설정
    multiline.negate: true
    multiline.match: after

output.logstash:
  hosts: ["localhost:5044"]

# 또는
output.elasticsearch:
  hosts: ["http://localhost:9200"]
  index: "app-logs-%{+yyyy.MM.dd}"
```