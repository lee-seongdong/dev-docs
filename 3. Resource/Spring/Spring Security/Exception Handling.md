# Exception Handling
## 1. ExceptionTranslationFilter
Spring Security의 보안 예외를 처리하는 핵심 `Filter`  
`SecurityFilterChain`에 연결되어, 이후 필터들에서 발생하는 예외 처리  
`AuthenticationException`과 `AccessDeniedException`을 처리하여 적절한 응답을 생성하는 역할  

### *ExceptionTranslationFilter psudocode*
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


## 2. AuthenticationException
인증이 실패했을 때 발생하는 예외  
사용자의 자격 증명이 유효하지 않거나, 인증 프로세스에서 문제가 발생했을 때 발생  
이때, `SecurityContextHolder`를 clear 하고, `AuthenticationEntryPoint`를 통해 인증 요청


## 3. AuthenticationEntryPoint
인증이 필요한 클라이언트에게 자격 증명을 요청하는 인터페이스  
클라이언트에게 인증 방법을 제시하는 역할  
로그인페이지 리다이렉트, HTTP Basic, Bearer Token, 403코드 반환등 여러 인증 방식을 지원함


## 4. AccessDeniedException
인증은 성공했지만, 해당 리소스에 대한 권한(Authorization)이 없을 때 발생하는 예외


## 5. AccessDeniedHandler
`AccessDeniedException`발생 시 호출되는 인터페이스  
일반적으로 `403 Forbidden` 응답을 반환하거나 에러페이지로 리다이렉트  


## 6. 예외 처리 흐름
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
