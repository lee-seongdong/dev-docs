## 이슈
Spring Bean이 아닌 클래스에서, Spring Bean을 사용해야 하는 이슈가 있음

## 해결
ApplicationContext을 통해, 원하는 곳에서 등록된 Bean을 사용
```java
public class SpringContextUtil implements ApplicationContextAware {
	private static ApplicationContext context;

	@Override
	public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
		context = applicationContext;
	}

	public static <T> T getBean(Class<T> beanClass) {
		return context.getBean(beanClass);
	}

	public static Object getBean(String beanName) {
		return context.getBean(beanName);
	}
}
```

```java
MyBean bean = SpringContextUtil.getBean(MyBean.class);
```