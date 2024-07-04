## 이슈
spring 3.x 에서 버전 호환성이슈로 인해 Redis Cluster Connection Factory 사용에 이슈가 있음

## 해결
apache commons-pool2를 사용해서 lettuce.RedisClusterConnection pool을 직접 구현

```java
public class RedisConnectionPool {
	private final RedisClusterClient clusterClient;
	private final GenericObjectPool<StatefulRedisClusterConnection<String, String>> pool;

	public RedisConnectionPool() {
		ClientResources clientResources = DefaultClientResources.create();

		RedisURI node1 = RedisURI.create("10.113.160.141", 22899);
		RedisURI node2 = RedisURI.create("10.113.160.141", 22909);

		clusterClient = RedisClusterClient.create(clientResources, Arrays.asList(node1, node2));

		GenericObjectPoolConfig<StatefulRedisClusterConnection<String, String>> poolConfig = new GenericObjectPoolConfig<>();
		poolConfig.setMaxTotal(8); // 최대 활성 연결 수
		poolConfig.setMaxIdle(8); // 최대 유휴 연결 수
		poolConfig.setMinIdle(0); // 최소 유휴 연결 수

		pool = ConnectionPoolSupport.createGenericObjectPool(clusterClient::connect, poolConfig);
	}

	public StatefulRedisClusterConnection<String, String> getConnection() throws Exception {
		return pool.borrowObject();
	}

	public void returnConnection(StatefulRedisClusterConnection<String, String> connection) {
		pool.returnObject(connection);
	}

	public void shutdown() {
		pool.close();
		clusterClient.shutdown();
	}

	public static void main(String[] args) {
		RedisConnectionPool connectionPool = new RedisConnectionPool();

		try (StatefulRedisClusterConnection<String, String> connection = connectionPool.getConnection()) {
			RedisAdvancedClusterCommands<String, String> commands = connection.sync();
			String value = commands.get("test");
			System.out.println("Stored value: " + value);
		} catch (Exception ignored) {

		} finally {
			connectionPool.shutdown();
		}
	}
}

```