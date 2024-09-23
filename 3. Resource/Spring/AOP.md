## Spring AOP
비즈니스 로직과 무관한 `횡단 관심사`를 모듈화 하여 프로그래밍 하는 기법이다.  
비즈니스 로직과 부가기능을 분리하여 가독성과 유지보수성을 높일 수 있고, 중복코드를 제거할 수 있다.  
관심사가 많아질수록 코드 흐름을 이해하기 어려워 가독성이 낮아질 수 있다.

## Spring AOP의 구현 방식
모듈화된 로직을 실행하기 위해서는 비즈니스 로직과의 결합이 필요한데, 이 과정을 `Weaving`이라고 한다.  
Spring에서는 `RTW` 방식으로 AOP를 구현했는데, 이유는 다음과 같다.
- 컴파일러와 JVM에 별도의 설정이 필요하지 않아 간단하게 사용할 수 있다.
- 특정 Bean이나 메소드에만 선택적 으로 AOP 적용이 가능함 (CTW, LTW 방식은 모든 대상클래스에 AOP가 적용됨)

> Weaving 종류  
> \- Compile Time Weaving (CTW) : 컴파일 과정에서 AOP 로직을 대상 객체에 삽입하는 방식. 성능면에서 가장 좋지만, lombok과 같이 컴파일 과정에서 코드를 조작하는 다른 플러그인과 충돌이 발생할 수 있다.  
> \- Load Time Weaving (LTW) : 바이트코드를 JVM에 로드할때 AOP 로직을 삽입하는 방식.  
> \- **Runtime Weaving (RTW)** : 런타임에 *프록시 객체*를 생성하여, 메소드 호출 시 AOP 로직을 적용하는 방식. 소스파일, 클래스파일의 변형이 필요하지 않아 별도의 설정이 필요없이 사용할 수 있지만, 관심사가 많아지면 오버헤드가 생겨 성능이 하락할 수 있다.

## Runtime Weaving과 프록시 객체
Runtime Weaving은 `JDK Dynamic Proxy`또는 `CGLib` 방식으로 구현할 수 있다.
### JDK Dynamic Proxy
Interface를 기반으로 프록시 객체를 생성하기 때문에, 반드시 Interface가 있어야 한다.  

Spring @Transactional 프록시 객체를 구현해보자.
#### 1. Interface와 Target 객체 구현
```java
public interface Calculable {
	int sum(int from, int to);
}
```
```java
public class TargetClass implements Calculable {
	@Override
	public int sum(int from, int to) {
		System.out.println("Call Target Method...");
		int sum = 0;
		for (int i = from; i <= to; i++) {
			sum += i;
		}
		return sum;
	}
}
```

#### 2. InvocationHandler를 implement하여 프록시 핸들러를 구현
```java
public class TransactionalHandler implements InvocationHandler {

	private final Object target;

	public TransactionalHandler(Object target) {
		this.target = target;
	}

	@Override
	public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
		try {
			System.out.println("start Transaction");
			Object result = method.invoke(target, args);
			System.out.println("end Transaction");

			return result;
		} catch (Exception e) {
			System.out.println("rollback Transaction");
			throw e;
		}
	}
}
```

#### 3. 구현한 프록시 핸들러를 통해 프록시 객체를 생성하고, 이를 통해 Target 객체 접근
```java
public class DynamicProxyTest {
	@Test
	public void testDynamicProxy() {
		Calculable target = new TargetClass();
		Calculable transactionalProxy = (Calculable)Proxy.newProxyInstance(
			target.getClass().getClassLoader(),
			target.getClass().getInterfaces(),
			new TransactionalHandler(target)
		);
		int ans = transactionalProxy.sum(1, 100);
		System.out.println("transactional Proxy ans : " + ans);
	}
}
```


### CGLib
Class를 기반으로 프록시 객체를 생성하기 때문에, Interface가 필요하지 않다.  

Spring @Async 프록시 객체를 구현해보자.
#### 1. Target 객체 구현
```java
public class TargetClass {
	public CompletableFuture<Integer> sumNumbers() throws InterruptedException {
		Thread.sleep(3000);

		System.out.println("Call Target Method...");
		int sum = 0;
		for (int i = 0; i < 10000; i++) {
			sum += i;
		}

		return CompletableFuture.completedFuture(sum);
	}

	public void printNumbers() throws InterruptedException {
		Thread.sleep(3000);

		System.out.println("Call Target Method...");
		for (int i = 0; i < 10000; i++) {
			System.out.println(i);
		}
	}
}
```

#### 2. MethodInterceptor 구현하여 비동기 로직 구현 (리턴 타입에 따른 분기 처리)
```java
public class AsyncInterceptor implements MethodInterceptor {
	private final ExecutorService executor = Executors.newFixedThreadPool(4);

	@Override
	public Object intercept(Object obj, Method method, Object[] args, MethodProxy proxy) {
		if (CompletableFuture.class.isAssignableFrom(method.getReturnType())) {
			return CompletableFuture.supplyAsync(() -> {
				System.out.println("start Async");
				try {
					Object result = proxy.invokeSuper(obj, args);
					return ((CompletableFuture<?>)result).join();
				} catch (Throwable e) {
					throw new RuntimeException(e);
				}
			}, executor);
		} else {
			CompletableFuture.runAsync(() -> {
				System.out.println("start Async");
				try {
					proxy.invokeSuper(obj, args);
				} catch (Throwable e) {
					throw new RuntimeException(e);
				}
			});
			return null;
		}
	}

	public void shutdown() {
		executor.shutdown();
	}
}
```

#### 3. Enhancer와 MethodInterceptor를 통해 프록시 객체를 생성하고, 이를 통해 Target 객체 조작
```java
public class CGLibProxyTest {
	@Test
	public void testCGLibProxyWithReturnValue() throws InterruptedException {
		TargetClass target = new TargetClass();
		AsyncInterceptor asyncInterceptor = new AsyncInterceptor();

		Enhancer enhancer = new Enhancer();
		enhancer.setSuperclass(target.getClass());
		enhancer.setCallback(asyncInterceptor);
		TargetClass proxy = target.getClass().cast(enhancer.create());

		CompletableFuture<Integer> futureResult = proxy.sumNumbers();
		futureResult.thenAccept(integer -> System.out.println("result: " + integer));

		asyncInterceptor.shutdown();
	}

    @Test
	public void testCGLibProxyWithReturnVoid() throws InterruptedException {
		TargetClass target = new TargetClass();
		AsyncInterceptor asyncInterceptor = new AsyncInterceptor();

		Enhancer enhancer = new Enhancer();
		enhancer.setSuperclass(target.getClass());
		enhancer.setCallback(asyncInterceptor);
		TargetClass proxy = target.getClass().cast(enhancer.create());

		proxy.printNumbers();

		Thread.sleep(5000);
		asyncInterceptor.shutdown();
	}
}
```


