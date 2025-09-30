# Spring Security Main Architecture
## 1. DelegatingFilterProxy
Servlet 컨테이너와 SpringContext 간의 연결을 담당하는 `Filter`  
Servlet 컨테이너가 Spring Bean을 인식할 수 있도록 위임(delegate)하는 역할


## 2. FilterChainProxy
Spring Security에서 Bean으로 제공하는, 보안 처리의 진입점 역할을 하는 `Filter`  
요청 URL패턴에 따라 적절한 `SecurityFilterChain`을 선택하여 실행하는 역할  
`DelegatingFilterProxy`를 통해 `FilterChain`에 연결되고, 실제 보안 필터링은 `SecurityFilterChain`에서 수행


## 3. SecurityFilterChain
보안 기능을 수행할 `SecurityFilter`들을 순서대로 정의한 컴포넌트  
하나의 애플리케이션에서 여러 개의 `SecurityFilterChain`을 구성하여 각기 다른 보안 정책을 적용할 수 있음  
커스텀 보안 필터는 `SecurityFilterChain`에 등록해야함  


## 4. SecurityFilter
Spring Security에서 제공하는 보안 `Filter`  
인증, 인가, CSRF 보호 등 다양한 보안 기능을 수행


## 5. 관계도
```
HTTP Request
    ↓
Servlet Container
    ↓
DelegatingFilterProxy (Servlet Filter)
    ↓ (위임)
FilterChainProxy (Spring Bean)
    ↓ (필터 체인 선택)
SecurityFilterChain{1, 2, 3, ...}
    ↓
[SecurityFilter1 → SecurityFilter2 → SecurityFilter3]
    ↓
DispatcherServlet
```