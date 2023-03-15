Cmocka Quick Start
==================

The project helps programmers by a quick start example.
The source files are from [cmocka-1.1.7 examples][cmocka-1.1.7 examples].

Just install cmocka and then run examples in this project.

Examples
--------

* check memory issues (mem leak, buffer overflow, etc.)
* replace subfunction with a mocked one ([uptime](./uptime))

Check Memory Issues
-------------------

This section shows how to write a unit test
(concerning memory leak, buffer overflow/underflow),
how to compile and run the testcases.

```c
# compile the code and generate executable file (the testcase exe)
> gcc -o allocate_module_test allocate_module_test.c allocate_module.c -lcmocka -DUNIT_TESTING=1

# run the test case
> ./allocate_module_test
output:
[==========] tests: Running 3 test(s).
[ RUN      ] leak_memory_test
[  ERROR   ] --- Blocks allocated...
allocate_module.c:41: note: block 0x5645fca1a750 allocated here
ERROR: leak_memory_test leaked 1 block(s)

[  FAILED  ] leak_memory_test
[ RUN      ] buffer_overflow_test
[  ERROR   ] --- allocate_module.c:48: error: Guard block of 0x5645fca1a7b0 size=4 is corrupt
allocate_module.c:46: note: allocated here at 0x5645fca1a7b4
[   LINE   ] --- allocate_module.c:48: error: Failure!
[  FAILED  ] buffer_overflow_test
[ RUN      ] buffer_underflow_test
[  ERROR   ] --- allocate_module.c:54: error: Guard block of 0x5645fca1a990 size=4 is corrupt
allocate_module.c:52: note: allocated here at 0x5645fca1a98f
[   LINE   ] --- allocate_module.c:54: error: Failure!
[  FAILED  ] buffer_underflow_test
[==========] tests: 3 test(s) run.
[  PASSED  ] 0 test(s).
[  FAILED  ] tests: 3 test(s), listed below:
[  FAILED  ] leak_memory_test
[  FAILED  ] buffer_overflow_test
[  FAILED  ] buffer_underflow_test

 3 FAILED TEST(S)
```

As we can see, the testcase reports memory leaks, buffer overflow/underflow.

[cmocka-1.1.7 examples]: https://git.cryptomilk.org/projects/cmocka.git/tree/example?h=cmocka-1.1.7
