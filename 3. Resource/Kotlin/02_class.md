- Visibility Modifier (= Access Modifier)
  - private : 동일 클래스 내에서만
  - public : 모든 파일에서 접근가능 (default)
  - protected : 상속한 자식 객체에서만
  - internal : 동일 모듈에서만
- Class
  - initial block : 주 생성자의 일부로 실행되는 코드블록. 불변 property의 값을 설정할 수 있다.
    ```kotlin
    class Person {
        val firstProperty = "First property".also(::println)
        init {
            println("First init block $firstProperty")
        }
        
        val secondProperty = "Second property".also(::println)
        init {
            println("Second init block $secondProperty")
        }
    }
    // First property
    // First init block First property
    // Second property
    // Second init block Second property
    ```
  - 주 생성자
    ```kotlin
    class Person constructor(firstName: String) { /*...*/ }
    class Person private constructor(firstName: String) { /*...*/ }
    class Person private @Inject constructor(firstName: String) { /*...*/ }

    // visibility modifier, annotation이 없는 constructor 는 생략가능
    class Person(name: String) { /*..*/ }

    // 주 생성자에서 전달받은 인자는 클래스 property나 initial block에 바로 사용 가능
    class Person(inputName: String) {
        val name = inputName
        init {
            println(inputName)
        }
    }

    // 객체의 property로 사용할 경우, val/var 키워드로 간결하게 사용가능
    class Person(val name: String) { /*..*/ }

    // default parameter
    class Person(val name: String, val student: Boolean = false) { /*..*/ }
    ```
  - 부 생성자
    ```kotlin
    class Person {
        val name
        // 클래스 내부에 constructor 키워드로 사용
        // val/var 키워드로 객체의 property로 바로 사용할 수 없음
        constructor(name: String) {
            this.name = name
        }
    }

    class Person(val name: String, val age: Int) {
        // 부 생성자는 this 키워드로 다른 생성자에게 객체 생성을 위임할 수 있다.
        constructor(name: String, age: Int, something: String) : this(name) { 
            /*..*/
        }
        
        // 주 생성자가 있다면, 객체 생성은 반드시 주 생성자에게 위임해야한다.
        constructor(name: String) : this(name, 0) {
            /*..*/
        }
    }
    ```
    - 생성자는 `주 생성자(init block) -> 부 생성자` 순으로 호출된다.
- 상속
  - 모든 클래스는 Any 를 상속한다 (= java의 Object)
  - 기본적으로 class는 final속성을 가진다. 상속을 혀용하려면 open 키워드가 필요하다.
    ```kotlin
    open class BaseClassA
    class DerivedClassA1 : BaseClassA()
    class DerivedClassA2(val x: Int) : BaseClassA()
    ```
- Properties
  - getter / setter : kotlin에서는 컴파일타임에 getter / setter를 생성
    - 선언 방식에 따라 생성 유무를 판단
      |클래스 선언|getter 생성|setter 생성||
      |---|---|---|---|
      |class Person(name: String)|X|X|주 생성자 매개변수|
      |class Person(var name: String)|O|O|속성|
      |class Person(val name: String)|O|X|속성|
    ```kotlin
    class Student(
        private val id: String, // private 불변 속성
        var name: String, // public 가변 속성
        age: Int) { // 생성자 매개변수

        // custom getter / setter
        // 어노테이션이나 visibility modifier를 설정할 수 있다.
        var age = age
            @Inject set(value) {
                require(value >= 0) { "Grade must be positive" }
                field = value
            }
            private get() {
                return field
            }

        // computed property
        val isAdult: Boolean
            get() {
                return age >= 18
            }
        
        // backing field
        var counter = 0
            set(value) {
                if (value >= 0) {
                    // counter = value 사용 시 무한 루프 (setter 재귀호출)
                    field = value // field 키워드로 setter 호출 방지 (backing field)
                }
            }
    }
    ```
- Extension
  - 기존 클래스를 수정하지 않고, 새로운 프로퍼티나 함수를 추가할 수 있는 기능. 특히 서드파티 라이브러리에 사용될때 유용하다.
  - Extension Function
    ```kotlin
    class Student(private val id: String, val name: String) {
    }

    fun Student.speak() {
        // println("Hello, my id is $id") // private 접근 불가
        println("Hello, my name is $name")
    }
    ```
  - Extension은 정적으로 결정되므로, 선언된 타입에 따라 확장함수나 프로퍼티를 호출한다.
    ```kotlin
    open class BaseClass
    class SubClass : BaseClass()

    fun BaseClass.printExtension() {
        println("BaseClass Extension")
    }

    fun SubClass.printExtension() {
        println("SubClass Extension")
    }

    val c = SubClass()
    c.printExtension() // SubClass Extension

    val c: BaseClass = SubClass()
    c.printExtension() // BaseClass Extension
    ```
  - Companion Object에도 적용이 가능하다.
    ```kotlin
    class BaseClass {
        companion object {}
    }

    fun BaseClass.Companion.print() {
        println("static extension fun")
    }

    BaseClass.Companion.print()
    BaseClass.print() // .Companion은 생략이 가능함
    ```