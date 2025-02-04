# Unit Test

Automatic testcases make life easer for developers.
After adding a new feature, the programmer can add
new cases to the pool of the existing testcases. And
then run (corresponding) testcases to keep code quality.

## C language

[cmocka][cmocka] is unit testing framework for C.  
[cmocka quick start][cmocka quick start] (in this repo)
just "show me the code".

## Golang

golang environment makes writing test cases very easy
via its standard package [testing][ut golang testing pkg].

[test case examples][ut golang examples] gives a quick start.  

## Rust

Rust provides a convenient way to write unit tests.
Just add testcases within the `tests` module.

[How to Write Tests][ut rust TRPL]

[cmocka]: https://cmocka.org/
[cmocka quick start]: ./cmocka
[ut golang examples]: https://github.com/hzget/go-investigation/tree/main/testing
[ut golang testing pkg]: https://pkg.go.dev/testing
[ut rust TRPL]: https://rust-book.cs.brown.edu/ch11-01-writing-tests.html
