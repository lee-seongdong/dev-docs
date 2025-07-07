# Spring Security
Spring 기반 애플리케이션에 인증(Authentication)과 인가(Authorization) 기능을 제공하는 보안 프레임워크


## Main Architecture
### 1. DelegatingFilterProxy
Servlet 컨테이너와 SpringContext 간의 연결을 담당하는 `Filter`  
Servlet 컨테이너가 Spring Bean을 인식할 수 있도록 위임(delegate)하는 역할

### 2. FilterChainProxy
Spring Security에서 Bean으로 제공하는, 보안 처리의 진입점 역할을 하는 `Filter`  
요청 URL패턴에 따라 적절한 `SecurityFilterChain`을 선택하여 실행하는 역할  
`DelegatingFilterProxy`를 통해 `FilterChain`에 연결되고, 실제 보안 필터링은 `SecurityFilterChain`에서 수행

### 3. SecurityFilterChain
보안 기능을 수행할 `SecurityFilter`들을 순서대로 정의한 컴포넌트  
하나의 애플리케이션에서 여러 개의 `SecurityFilterChain`을 구성하여 각기 다른 보안 정책을 적용할 수 있음  
커스텀 보안 필터는 `SecurityFilterChain`에 등록해야함  

### 4. SecurityFilter
Spring Security에서 제공하는 보안 `Filter`  
인증, 인가, CSRF 보호 등 다양한 보안 기능을 수행

### 관계도
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


## Authentication Architecture
### 1. SecurityContextHolder
인증된 사용자의 세부 정보인 `SecurityContext`를 저장하는 핵심 컴포넌트  
현재 인증된 사용자 정보에 접근할 수 있는 진입점 역할  
기본적으로 `ThreadLocal`을 사용하여 요청 스레드별로 독립적으로 동작

### 2. SecurityContext
인증된 사용자의 세부 정보를 담고있는 컴포넌트  
현재 인증된 사용자의 `Authentication` 객체를 포함  

### 3. Authentication
Spring Security에서 인증을 위한 핵심 인터페이스  
세 가지 요소로 구성됨:
- Principal: 식별된 사용자 객체. 대개 `UserDetails`를 구현한 객체
- credentials: 사용자 인증에 사용할 비밀번호. 일반적으로 사용자 인증 후 유출되지 않게 clear함
- authorities: 사용자에게 부여된 권한. `GrantedAuthority`를 구현한 객체

두 가지 주요 용도로 사용됨:
- `AuthenticationManager`의 입력으로 사용되어, 사용자가 제공한 자격 증명을 전달
- `SecurityContext`에서 현재 인증된 사용자를 나타냄

#### 인증 객체의 관계도
```
SecurityContextHolder
    ↓ (포함)
SecurityContext
    ↓ (포함)
Authentication
    ↓ (포함)
[Principal, Credentials, GrantedAuthority]
```

### 4. UserDetails
사용자 정보를 제공하는 인터페이스  
사용자명, 비밀번호, 권한, 계정 상태 등을 포함하는 사용자 정보의 표준 형태  
Spring Security에서 직접 보안 목적으로 사용하지 않고, 사용자 정보를 저장하여 `Authentication` 객체로 캡슐화됨  

### 5. GrantedAuthority
사용자에게 부여된 권한(`ROLE_ADMIN`, `ROLE_USER`)을 나타내는 인터페이스  
일반적으로 역할(roles)이나 스코프(scopes)를 의미함  

### 6. AuthenticationManager
`SecurityFilter`들이 인증을 수행하는 방법을 정의한 인터페이스  
인증이 성공하면 반환된 `Authentication`이 `SecurityContextHolder`에 저장됨  
인증정보를 직접 `SecurityContextHolder`에 저장하는 경우, `AuthenticationManager`는 사용할 필요 없음

### 7. ProviderManager
`AuthenticationManager`의 가장 일반적인 구현체  
`AuthenticationProvider`들을 리스트형태로 가지며, 요청받은 인증 유형에 해당하는 Provider에게 인증을 위임함  

모든 Provider에서 실패하는 경우, `ProviderNotFoundException`와 함께 인증 실패

### 8. AuthenticationProvider
특정 유형의 인증을 수행하는 인터페이스  
예를 들어 `DaoAuthenticationProvider`는 사용자명/비밀번호 기반 인증을, `JwtAuthenticationProvider`는 JWT 기반 인증을 수행

### 9. UserDetailsService
`DaoAuthenticationProvider`에서 사용자 정보를 로드하기 위한 인터페이스  
```
UserDetailsService.loadUserByUsername("user123")
    ↓
UserDetails 객체 반환
    ↓
DaoAuthenticationProvider가 사용
    ↓
Authentication 객체 생성
    ↓
SecurityContextHolder에 저장
```

### 10. AbstractAuthenticationProcessingFilter
요청으로부터 `Authentication`을 생성하기 위한 추상 `Filter`  
생성된 `Authentication`은 `AuthenticationManager`을 통해 인증 처리  
구현체로는 `UsernamePasswordAuthenticationFilter`, `OneTimeTokenAuthenticationFilter`등이 있음

### 인증 처리 흐름
```
HTTP Request
    ↓
AbstractAuthenticationProcessingFilter
    ↓
1. HttpServletRequest로부터 Authentication(미인증) 생성
    ↓
2. AuthenticationManager(ProviderManager)로 Authentication 전달하여 인증 요청
    ↓
3. AuthenticationManager에서 적절한 AuthenticationProvider 선택
    ↓
4. AuthenticationProvider가 실제 인증을 위임받아 수행
    ↓
<인증 결과 확인>
    1. 인증 성공
    → SecurityContextHolder에 Authentication(인증됨) 저장
    → AuthenticationSuccessHandler 호출

    2. 인증 실패
    → SecurityContextHolder 클리어
    → AuthenticationFailureHandler 호출
```


## Exception Handling
### 1. ExceptionTranslationFilter
Spring Security의 보안 예외를 처리하는 핵심 `Filter`  
`SecurityFilterChain`에 연결되어, 이후 필터들에서 발생하는 예외 처리
`AuthenticationException`과 `AccessDeniedException`을 처리하여 적절한 응답을 생성하는 역할  

#### *ExceptionTranslationFilter psudocode*
```java
try {
	filterChain.doFilter(request, response); 
} catch (AuthenticationException | AccessDeniedException ex) {
	if (!authenticated || ex instanceof AuthenticationException) {
		startAuthentication(); 
	} else {
		accessDenied(); 
	}
}
```

### 2. AuthenticationException
인증이 실패했을 때 발생하는 예외  
사용자의 자격 증명이 유효하지 않거나, 인증 프로세스에서 문제가 발생했을 때 발생  
이때, `SecurityContextHolder`를 clear 하고, `AuthenticationEntryPoint`를 통해 인증 요청

### 3. AuthenticationEntryPoint
인증이 필요한 클라이언트에게 자격 증명을 요청하는 인터페이스  
클라이언트에게 인증 방법을 제시하는 역할  
로그인페이지 리다이렉트, HTTP Basic, Bearer Token, 403코드 반환등 여러 인증 방식을 지원함

### 4. AccessDeniedException
인증은 성공했지만, 해당 리소스에 대한 권한(Authorization)이 없을 때 발생하는 예외

### 5. AccessDeniedHandler
`AccessDeniedException`발생 시 호출되는 인터페이스  
일반적으로 `403 Forbidden` 응답을 반환하거나 에러페이지로 리다이렉트  

### 예외 처리 흐름
```
SecurityFilterChain
    ↓
[SecurityFilter1 → ExceptionTranslationFilter → SecurityFilter2 → SecurityFilter3]
    ↓
ExceptionTranslationFilter 다음 필터에서 예외 발생 (AccessDeniedException | AuthenticationException)
    ↓
<사용자 인증 상태 확인>
    1. 미인증 OR AuthenticationException 
    → startAuthentication() 
    → AuthenticationEntryPoint.commence() 
    → 로그인 페이지 리다이렉트 또는 WWW-Authenticate 헤더 전송

    2. 인증됨 AND AccessDeniedException 
    → accessDenied() 
    → 에러페이지 리다이렉트 또는 403 응답
```

