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
- Extension Lambda
  - 
- scope function
  - 주어진 객체에 대한 scope function을 호출하면, 임시 블록을 생성하고, 여기에서는 주어진 객체를 이름없이 접근 가능하다.
  - 이때, scope function을 호출한 객체를 수신객체라 한다.
  - scope function 종류 : apply, also, run, let, with 등
    ```kotlin
    // 1. apply : 수신객체를 this로 사용하고, 수신객체를 반환
    // 객체의 초기화 및 설정에 적합
    val person = Person().apply {
        name = "Lee"
        age = 20
    }

    // 2. also : 수신객체를 파라미터로 전달하여 사용하고, 수신객체를 반환
    // 객체를 수정하지않는 부가작업에 적합
    val text = "Hello".also { print(it) } // text = Hello

    // 3. run : 수신객체를 this로 사용하고, 람다 결과값을 반환
    // 객체와 관련된 작업을 수행한 후 값을 반환할 때 적합
    val graterThanZero = 10.run { this > 0 }

    // 4. let : 수신객체를 파라미터로 전달하여 사용하고, 람다 결과값을 반환
    // 수신객체를 변환하거나, null check에 적합
    val length = "Hello".let { it.length }

    // 5. with : 수신객체를 this로 사용하고, 람다의 결과값을 반환
    // 수신객체를 컨텍스트로 작업을 수행할때 적합
    val graterThanZero = with(10) { this > 0 }
    ```