__작성중
## CCR(Cross Cluster Replication)
### 1. 팔로워 클러스터에서 Remote Cluster Connection 설정
```bash
# API를 통해 리더 클러스터 연결 등록 (팔로워 클러스터에서 실행)
PUT /_cluster/settings
{
  "persistent": {
    "cluster.remote": {
      "my-leader-cluster": { # 이 이름이 leader_alias가 됨
        "seeds": ["leader-node1:9300", "leader-node2:9300"] 
      }
    }
  }
}
```

### 2. 복제 시작 (팔로워 클러스터에서 실행)
```bash
# 특정 인덱스 복제 시작
PUT /_plugins/_replication/follower-index/_start
{
  "leader_alias": "my-leader-cluster",         # 위에서 설정한 cluster.remote 이름
  "leader_index": "source-index",
  "use_roles": {
    "leader_cluster_role": "all_access",
    "follower_cluster_role": "all_access"
  }
}
```

### 3. 자동 팔로우 패턴 설정 (선택사항)
```bash
# 패턴 기반 자동 복제 규칙
POST /_plugins/_replication/_autofollow
{
  "leader_alias": "my-leader-cluster",         # cluster.remote에서 설정한 이름
  "name": "logs-auto-follow",
  "pattern": "logs-*",
  "use_roles": {
    "leader_cluster_role": "all_access",
    "follower_cluster_role": "all_access"
  }
}
```

### 4. 복제 상태 확인
```bash
# 복제 상태 확인
GET /_plugins/_replication/follower-index/_status

# 자동 팔로우 통계 확인
GET /_plugins/_replication/autofollow_stats
```