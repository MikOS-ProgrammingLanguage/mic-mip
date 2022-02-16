# Assignments in LLVM-IR

## Assign directly

To assign a value like ``int test_int = 100`` you would need to do the following

```c
define i32 @main() {
    %test_int = alloca i32*, align 4
    store i32 100 = i32* %test_int, align 4
}
```

## Assign expression

To assign a expression like ``int test_int = 100 + 20 * 30`` you would need to do the following

```c
define i32 @main() {
    %test_int = alloca i32*, align 4 // allocate a i32*. Align 4 means, that the memory address is a multiple of 4
    %1 = (type prefix)mul i32 20, 30 // multiplies 20 and 30
    %2 = (type prefix)add i32 %1, 100 // adds %1 (= 20*30) and 100
    store i32 %2, i32* %test_int, align 4 // stores %2 (= 100 + 20 * 30) into %test_int
}
```

## Assign expression with references

To assign a expression with a reference like ``int test_int = a`` you would need to do the following

```c
define i32 @main() {
    %a = alloca i32*, align 4
    store i32 10, i32* %a
    %test_int = alloca i32*, align 4
    %1 = load i32, i32* %a, align 4
    store i32 %1, i32* %test_int
}
```
