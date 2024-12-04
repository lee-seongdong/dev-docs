- 프로그램 진입점은 main 함수
  ```kotlin
  fun main() {
      println("hello, world")
  }
  ```
- 상수 val (불변)
  ```kotlin
  // 선언 후 초기화
  val a: Int
  a = 10

  // 선언 및 초기화
  val b: Int = 11

  // 선언 및 초기화(타입추론)
  val c = 12
  ```
- 변수 var (가변)
  ```kotlin
  // 선언 후 초기화
  var a: Int
  a = 10

  // 선언 및 초기화
  var b: Int = 11

  // 선언 및 초기화(타입추론)
  var c = 12
  ```
- 함수 fun : 이름이 붙은 서브루틴
    ```kotlin
    fun funA(arg1: Int, arg2: Type, ...): ReturnType {
        var result: ReturnType = "test"
        return result
    }

    // 식 본문 함수
    fun plus(a: Int, b: Int) = a + b
    ```
- if : if도 식이기 때문에, 값을 반환한다.  
  ```kotlin
  val a = if (true) 10 else 20

  // 블록으로도 가능함. 단, return 키워드는 제외
  val b = if (true) {
      println("true")
      10
  } else {
      println("false")
      20
  }
  ```
- 문자열 템플릿
  ```kotlin
  val number = 10
  println("Number is $number")
  println("${if (true) 'a' else 'b'}")
  
  val name = "John"
    val age = 30
    val json = """{
    "name": "$name",
    "age": $age
  }"""
    println(json)
  ```
- for
  ```kotlin
  for (i in 1..5) print("$i ") // 1 2 3 4 5
  for (i in 1 until 5) print("$i ") // 1 2 3 4
  for (i in 10 downTo 1 step 3) print("$i ") // 10 7 4 1
  for (i in 1..10 step 2) // 1 3 5 7 9
  for (c in 'A'..'Z' step 2) // A C E G I K M O Q S U W Y
  ```