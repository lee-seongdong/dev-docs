# Command

## 리다이렉션
`>` : stdout 리다이렉션(덮어쓰기)
`>>` : stdout 리다이렉션(이어쓰기)
`2>` : stderr 리다이렉션
`2>&1` : stderr를 1번 디스크립터(stdout)로 보내기
`2>1` : stderr를 1번 파일에 쓰기
`<` : stdin 리다이렉션 (파일을 입력으로)
`|` : 파이프 (stdout -> stdin)

## 파일 디스크립터
모든 프로세스는 3개의 기본 스트림을 가짐:

0 = stdin  (표준 입력)  ← 키보드 입력
1 = stdout (표준 출력)  → 화면 출력
2 = stderr (표준 에러)  → 에러 메시지
3~ = 파일 디스크립터(FD) → 열린 파일을 가리키는 숫자

## User & Group
```
사용자 (User)
    ↓ 속함
그룹 (Group)
    ↓ 소유
파일/디렉터리
```

### User
```bash
# user = 시스템에서 작업하는 주체
# 사용자 정보 저장 위치
/etc/passwd   # 사용자 목록
/etc/shadow   # 비밀번호 (암호화)

# /etc/passwd 형식
irteam:x:10001:10001::/home1/irteam:/bin/bash
  │    │   │     │    │      │           │
  이름  │  UID   GID  설명   홈         쉘
      비밀번호(x=shadow에 있음)

# 주요 명령어
useradd -u 10001 -g irteam irteam   # 사용자 생성
userdel -r irteam                   # 사용자 삭제
id irteam                           # 정보 확인
whoami                              # 현재 사용자
```

### Group
```bash
# group = 사용자들의 모임 (권한 공유)
# 그룹 정보 저장 위치
/etc/group

# /etc/group 형식
irteam:x:10001:alice,bob
  │    │   │      │
 이름  │  GID   멤버 목록
    비밀번호(거의 안 씀)

# 주요 명령어
groupadd -g 10001 irteam     # 그룹 생성
groupdel irteam              # 그룹 삭제
groups irteam                # 그룹 목록 확인
```

0         = root (슈퍼유저)
1-999     = 시스템 사용자 (daemon, bin, ...)
1000+     = 일반 사용자

실제로는 숫자로 관리:
ls -la
drwxr-xr-x 2 10001 10001 4096 Jan 10 ...
              ^^^^^ ^^^^^
              UID   GID

"irteam"은 단지 10001의 별명!