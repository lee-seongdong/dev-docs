## 이슈
lombok의 @RequiredArgsConstructor 기본 설정으로 사용 시 몇가지 이슈가있다.
- @Value를 불변으로 사용하지 못함
- @Qualifier가 의도대로 동작하지 않음
```java
@Component
@RequiredArgsConstructor
public class MyClass {
    @Qualifier("sampleClient")
    private final WebClient webclient; // @Qualifier 애노테이션이 롬복이 생성한 생성자에 붙지않아 원하는 빈을 주입할 수 없다.

    @Value
    private final String apiUrl; // @Value 애노테이션이 롬복이 생성한 생성자에 붙지않아 값을 주입할 수 없다.
}
```

## 해결
lombok.copyableAnnotations 설정을 사용하여 해결  
위 설정은 프로퍼티의 애노테이션을 생성자로 복사할 애노테이션을 등록하는 설정이다.


```config
config.stopBubbling = true
lombok.copyableAnnotations += org.springframework.beans.factory.annotation.Value
lombok.copyableAnnotations += org.springframework.beans.factory.annotation.Qualifier
```

lombok이 생성하는 코드
```java
// copyableAnnotations 설정 이전
@Component
public class MyClass {
    @Qualifier("sampleClient")
    private final WebClient webclient;

    @Value
    private final String apiUrl;

    public MyClass(WebClient webclient, String apiUrl) {
        this.webclient = webclient;
        this.apiUrl = apiUrl;
    }
}
```

```java
// copyableAnnotations 설정 이후
@Component
public class MyClass {
    @Qualifier("sampleClient")
    private final WebClient webclient;

    @Value
    private final String apiUrl;

    public MyClass(@Qualifier("sampleClient") WebClient webclient, @Value String apiUrl) {
        this.webclient = webclient;
        this.apiUrl = apiUrl;
    }
}
```