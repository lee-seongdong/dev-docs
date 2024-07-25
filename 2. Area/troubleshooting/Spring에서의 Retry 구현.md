## Spring Retry
retry 정책은 back를 포함하는 것이 좋다.
```java
@Slf4j
@Component
public class RetryLoggingListener extends RetryListenerSupport {
	@Override
	public <T, E extends Throwable> void onError(RetryContext context, RetryCallback<T, E> callback, Throwable throwable) {
		log.warn("retry : {}..", context.getRetryCount());
		super.onError(context, callback, throwable);
	}

	@Override
	public <T, E extends Throwable> void close(RetryContext context, RetryCallback<T, E> callback, Throwable throwable) {
		log.error("retry close", throwable);
		super.close(context, callback, throwable);
	}
}
```

```java
@Configuration
public class RetryConfiguration {
	private static final int MAX_RETRIES = 3;
	private static final int INITIAL_RETRY_INTERVAL = 5 * 1000;
	private static final int MAX_RETRY_INTERVAL = 60 * 1000;
	private static final int MULTIPLIER = 2;

	@Bean
	public RetryTemplate retryTemplate(RetryLoggingListener retryLoggingListener) {
		RetryTemplate retryTemplate = new RetryTemplate();

		SimpleRetryPolicy retryPolicy = new SimpleRetryPolicy();
		retryPolicy.setMaxAttempts(MAX_RETRIES);
		retryTemplate.setRetryPolicy(retryPolicy);

		ExponentialBackOffPolicy backOffPolicy = new ExponentialBackOffPolicy();
		backOffPolicy.setInitialInterval(INITIAL_RETRY_INTERVAL);
		backOffPolicy.setMultiplier(MULTIPLIER);
		backOffPolicy.setMaxInterval(MAX_RETRY_INTERVAL);
		retryTemplate.setBackOffPolicy(backOffPolicy);

		retryTemplate.registerListener(retryLoggingListener);
		return retryTemplate;
	}
}
```
```java
...
    public String before() {
        String result = callApi();
        return result;
    }

    public String after() {
        String result = retryTemplate.execute(retryContext -> {
            String result = callApi();
            return result;
        });

        return result;
    }
...

```

<br/>

## WebClient Retry
```java
...
	webClient.get()
      .uri("localhost/test")
      .retrieve()
      .bodyToMono(String.class)
      .retryWhen(Retry.backoff(3, Duration.ofSeconds(2)));
...
```