- fun : 함수
  ```kotlin
  fun funA(arg1: Int, arg2: Type, ...): ReturnType {
      var result: ReturnType = "test"
      return result
  }

  // 식 본문 함수 : 식 자체를 본문으로 가지는 함수
  fun plus(a: Int, b: Int) = a + b

  // default parameter
  fun minus(a: Int = 100, b: Int) = a - b

  // named argument
  minus(b = 10)
  ```

  - kotlin은 함수가 1급객체이므로, 변수에 할당 및 반환이 가능
    ```kotlin
    // 1. 기본 함수
    fun minus(a: Int, b: Int): Int {
        return a - b
    }

    // 2. 익명함수를 변수에 할당
    val multiple = fun(a: Int, b: Int): Int {
        return a * b
    }

    // 3. lamda 표현식으로 변수에 함수 할당
    val plus = { a: Int, b: Int -> a + b }

    // 함수를 파라미터로 받는 함수
    fun operate(a: Int, b: Int, operation: (Int, Int) -> Int): Int {
        return operation(a, b)
    }

    operate(1, 2, plus)
    operate(1, 2, ::minus) // 기본 함수를 변수로서 전달하는 방법
    ```
- Functional Interface
  - 메소드를 하나만 가진 인터페이스를 functional interface  
  혹은 SAM(Single Abstract Method) interface 라고한다.
    ```kotlin
    // fun interface 키워드로 함수형 인터페이스 선언
    fun interface IntPredicate {
        fun accept(i: Int): Boolean
    }

    // SAM 변환을 사용하지 않은 경우
    val isEven = object : IntPredicate {
        override fun accept(i: Int): Boolean {
            return i % 2 == 0
        }
    }

    // SAM 변환을 사용하는 경우
    val isEven = IntPredicate { i -> i % 2 == 0 }
    ```
- lambda
  ```kotlin
  // 기본적인 forEach
  listOf(1, 2, 3, 4, 5).forEach(fun(value) {
      print(value)
  })
  // 함수참조를 사용한 forEach
  listOf(1, 2, 3, 4, 5).forEach(::print)

  // 1. 익명함수를 람다로 변환
  listOf(1, 2, 3, 4, 5).forEach({ value -> print(value) })

  // 2. 함수타입의 파라미터가 순서상 마지막에 오는 경우, 소괄호 밖으로 뺄 수 있음
  listOf(1, 2, 3, 4, 5).forEach() { value -> print(value) }
  // 이때, 나머지 파라미터가 없는 경우, 소괄호는 생략가능
  listOf(1, 2, 3, 4, 5).forEach { value -> print(value) }

  // 3. 함수타입의 파라미터가 하나의 매개변수만 사용하는 경우, it 키워드로 대체가능
  listOf(1, 2, 3, 4, 5).forEach { print(it) }
  ```