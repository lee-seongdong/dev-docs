# Spring Batch의 JobRepository란?
`JobRepository`는 Spring Batch에서 배치 작업의 메타데이터를 저장하고 관리하는 데 사용되는 인터페이스이다.  

## 배치 메타데이터
배치 메타데이터는 배치 작업의 상태와 진행 상황을 추적하는 데 사용되는 데이터이다.  
여기에는 작업 인스턴스, 작업 실행, 단계 실행 등에 대한 정보가 포함된다.  
주로 데이터베이스에 저장되며, 이를 통해 배치 작업의 상태를 추적하고 재시작 기능을 지원한다.

### 배치 메타데이터 테이블 종류
Spring Batch는 다음과 같은 기본 메타데이터 테이블을 사용한다.  
(org/springframework/batch/core/schema-{platform}.sql에서 확인할 수 있다.)
- `BATCH_JOB_INSTANCE`: 작업 인스턴스 정보
- `BATCH_JOB_EXECUTION`: 작업 실행 정보
- `BATCH_JOB_EXECUTION_PARAMS`: 작업 실행 파라미터
- `BATCH_STEP_EXECUTION`: 단계 실행 정보
- `BATCH_STEP_EXECUTION_CONTEXT`: 단계 실행 컨텍스트
- `BATCH_JOB_EXECUTION_CONTEXT`: 작업 실행 컨텍스트

### 테이블 생성
Spring Boot에서는 `spring.batch.initialize-schema` 설정으로 메타데이터 스키마 초기화에 대한 설정을 할 수 있다.
- always: 매번 초기화
- embedded: 내장 DB를 사용중일때만 초기화
- never: 초기화 안 함 (prod 환경에서 권장)

#### 테이블명 변경
기본 메타데이터 테이블명은 `JobRepositoryFactoryBean`의 `setTablePrefix`나 `setTableNames`로 변경할 수 있다.
- `setTablePrefix`
    ```java
    @Bean
    public JobRepository jobRepository(DataSource dataSource, PlatformTransactionManager transactionManager) throws Exception {
        JobRepositoryFactoryBean factory = new JobRepositoryFactoryBean();
        factory.setDataSource(dataSource);
        factory.setTransactionManager(transactionManager);
        factory.setDatabaseType(DatabaseType.MYSQL.name());
        factory.setTablePrefix("MY_BATCH_"); // 테이블 접두사 변경
        return factory.getObject();
    }
    ```
- `setTableNames`
    ```java
    @Bean
    public JobRepository jobRepository(DataSource dataSource, PlatformTransactionManager transactionManager) throws Exception {
        JobRepositoryFactoryBean factory = new JobRepositoryFactoryBean();
        factory.setDataSource(dataSource);
        factory.setTransactionManager(transactionManager);
        factory.setDatabaseType(DatabaseType.MYSQL.name());
        factory.setTableNames(new HashMap<String, String>() {{
            put("BATCH_JOB_INSTANCE", "MY_JOB_INSTANCE");
            put("BATCH_JOB_EXECUTION", "MY_JOB_EXECUTION");
            put("BATCH_JOB_EXECUTION_PARAMS", "MY_JOB_EXECUTION_PARAMS");
            put("BATCH_STEP_EXECUTION", "MY_STEP_EXECUTION");
            put("BATCH_STEP_EXECUTION_CONTEXT", "MY_STEP_EXECUTION_CONTEXT");
            put("BATCH_JOB_EXECUTION_CONTEXT", "MY_JOB_EXECUTION_CONTEXT");
        }}); // 테이블명 변경
        return factory.getObject();
    }
    ```
