- 프로그램 진입점은 main 함수
  ```kotlin
  fun main() {
      println("hello, world")
  }
  ```
- val (불변)
  ```kotlin
  // 선언 후 초기화
  val a: Int
  a = 10

  // 선언 및 초기화
  val b: Int = 11

  // 선언 및 초기화(타입추론)
  val c = 12
  ```
- var (가변)
  ```kotlin
  // 선언 후 초기화
  var a: Int
  a = 10

  // 선언 및 초기화
  var b: Int = 11

  // 선언 및 초기화(타입추론)
  var c = 12
  ```
- const val (상수)
  ```kotlin
  const val SYSTEM_NAME: String = "MyApp"
  ```
- lateinit
  ```kotlin
  lateinit var text: String // var 에만 사용 가능하고, Primitive Type에는 사용할 수 없다
  print(text) // Error!

  text = "text1"
  print(text) // text1

  text = "text2"
  print(text) // text2
  ```
- by lazy
  ```kotlin
  lateinit var text: String
  val textLength by lazy { // val 에만 사용 가능하다.
    text.length
  }

  print(textLength) // Error!

  text = "test"
  print(textLength) // 4

  text = "testtest"
  print(textLength) // 4
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

  // with index
  for (i in array.indices) print(array[i])
  for ((index, value) in array.withIndex()) println("the element at $index is $value")
  ```
- break to label
  ```kotlin
  loop@ for (i in 1..100) {
    for (j in 1..100) {
        if (...) break@loop
    }
  }
  ```
- return to label
  ```kotlin
  listOf(1, 2, 3, 4, 5).forEach {
      if (it == 3) return
      print(it)
  }
  // 12

  listOf(1, 2, 3, 4, 5).forEach {
      if (it == 3) return@forEach
      print(it)
  }
  // 1245

  listOf(1, 2, 3, 4, 5).forEach aa@{
      if (it == 3) return@aa
      print(it)
  }
  // 1245

  listOf(1, 2, 3, 4, 5).forEach(fun(value) {
      if (value == 3) return
      print(value)
  })
  // 1245
  ```
- when : 구문과 표현식 모두 사용가능. switch와 유사함
  ```kotlin
  // 표현식
  val text = when (x) {
      is String -> "x is string"
      1 -> "x == 1"
      2 -> "x == 2"
      in 10..20 -> "x is 1*"
      else -> "x is neither 1 nor 2"
  }

  // 구문
  when (x) {
      is String -> print("x is string")
      1 -> print("x == 1")
      2 -> print("x == 2")
      in 10..20 -> print("x is 1*")
      else -> print("x is neither 1 nor 2")
  }
  ```
- exception with precondition function
  |Precondition function|Use case|Exception thrown|
  |---|---|---|
  |require()|사용자 인풋을 검증|IllegalArgumentException|
  |check()|객체나 변수의 상태를 검증|IllegalStateException|
  |error()|비정상적인 상황에서 호출|IllegalStateException|
  - require
    ```kotlin
    require(str != null) // null인 경우 IllegalArgumentException 발생
    println(str.length) // require() 통과 시 String타입으로 smart cast
    ```
  - check
    ```kotlin
    check(str != null) // null인 경우 IllegalStateException 발생
    println(str.length) // check() 통과 시 String타입으로 smart cast
    ```
  - error
    ```kotlin
    if (str == null) error("str is null") // error()로 직접 예외 발생
    ```
destructuring : componentN() 연산자가 구현되어있는 객체에 한하여, 여러개의 값으로 분해할 수 있는 기능
  ```kotlin
  class Person {
      var name: String = ""
      var age: Int = 0

      // componentN 연산자
      operator fun component1(): String {
          return name
      }
      
      // componentN 연산자  
      operator fun component2(): Int {
          return age
      }
  }

  val (name, age) = Person().apply { name = "Alice"; age = 30 } // 구조분해
  ```