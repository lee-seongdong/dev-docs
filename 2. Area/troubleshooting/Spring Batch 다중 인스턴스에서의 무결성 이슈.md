## 이슈
1. Spring Batch 메타데이터를 H2DB로 저장하고 있음.
2. 배치 작업을 @Scheduled로 실행하고 있음.

**현 상태로 배치 인스턴스 증설 시, 데이터 무결성 이슈가 생길 수 있음.**

## 해결
1. 메타데이터를 H2DB가 아닌 외부 DB서버에 저장
    - H2DB는 인메모리 데이터베이스로, 다중 인스턴스 환경에서 데이터 무결성을 보장할 수 없다.  
    따라서, 외부 DB서버에 메타데이터를 저장하는 것이 좋다.
    - 예제 설정:
    ```java
    @Bean
    // JobRepository : 메타데이터를 관리하는데 사용됨
    public JobRepository jobRepository(DataSource dataSource, PlatformTransactionManager transactionManager) throws Exception {
        JobRepositoryFactoryBean factory = new JobRepositoryFactoryBean();
        factory.setDataSource(dataSource);
        factory.setTransactionManager(transactionManager);
        factory.setDatabaseType(DatabaseType.MYSQL.name());
        return factory.getObject();
    }
    ```

2. 배치 작업 실행 시, lock으로 중복 실행을 방지 (shedlock)
    - Lock을 사용하여 배치 작업의 중복 실행을 방지할 수 있다.  
    ShedLock은 분산 환경에서 동시 실행을 방지하기 위해 락을 제공한다.
    - ShedLock 설정 예제:
    ```java
    @Configuration
    @EnableScheduling
    @EnableSchedulerLock(defaultLockAtMostFor = "5m", defaultLockAtLeastFor = "5m")
    public class SchedulerConfig {
        
        @Bean
        public LockProvider lockProvider(DataSource dataSource) {
            return new JdbcTemplateLockProvider(
                JdbcTemplateLockProvider.Configuration.builder()
                    .withJdbcTemplate(new JdbcTemplate(dataSource))
                    .withTableName("BATCH_SHEDLOCK") // default : shedlock
                    .usingDbTime() // 배치 인스턴스의 시간이 아닌, DB의 시간을 사용하도록 설정
                    .build()
            );
        }
    }
    ```

    - 배치 작업에 ShedLock 적용 예제:
    ```java
    @Scheduled(cron = "0 0 * * * ?")
    @SchedulerLock(name = "scheduledTaskName")
    public void performTask() {
        // do something
    }
    ```

    - 참고 : [ShedLock GitHub](https://github.com/lukas-krecan/ShedLock/blob/master/README.md)

