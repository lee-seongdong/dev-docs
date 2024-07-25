## 이슈
외부 API의 rate limit 정책에 따라 배치 요청량을 제한해야 하는 이슈가 있음

## 해결
ChunkListener를 통해 배치의 각 chunk를 1분단위로 실행하도록 구현

```java
@Component
public class DelayedChunkListener extends ChunkListenerSupport {
	private static final long DELAY = 5 * 1000;

	@Override
	public void afterChunk(ChunkContext chunkContext) {
		try {
			Thread.sleep(DELAY);
		} catch (InterruptedException e) {
			Thread.currentThread().interrupt();
			throw new RuntimeException("Chunk interrupted", e);
		}
	}
}
```

```java
...

public Step testStep(DelayedChunkListener chunkListener) {
    return stepBuilderFactory.get("testStep")
			.<String, String>chunk(PAGE_SIZE)
			.reader(reader())
			.processor(processor())
			.writer(writer())
			.listener(chunkListener)
			.build();
}
...
```