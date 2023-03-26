Cmocka Quick Start
==================

The project helps programmers by a quick start example.
Some source files are modified from [cmocka-1.1.7 examples][cmocka-1.1.7 examples].

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

    // our lib to run a unit test
    // allocate_module.c

    #ifdef UNIT_TESTING
    extern void* _test_malloc(const size_t size, const char* file, const int line);
    extern void _test_free(void* const ptr, const char* file, const int line);

    #define malloc(size) _test_malloc(size, __FILE__, __LINE__)
    #define free(ptr) _test_free(ptr, __FILE__, __LINE__)
    #endif // UNIT_TESTING

    void leak_memory(void) {
        int * const temporary = (int*)malloc(sizeof(int));
        *temporary = 0;
    }

    void no_leak_memory(void) {
        int * const temporary = (int*)malloc(sizeof(int));
        free(temporary);
    }

    // add test cases in testing file
    // allocate_module_test.c

    extern void leak_memory(void);
    extern void no_leak_memory(void);

    static void leak_memory_test(void **state) {
        leak_memory();
    }
    static void no_leak_memory_test(void **state) {
        no_leak_memory();
    }

    int main(void) {
        const struct CMUnitTest tests[] = {
            cmocka_unit_test(leak_memory_test),
            cmocka_unit_test(no_leak_memory_test),
        };
        return cmocka_run_group_tests(tests, NULL, NULL);
    }

    # compile the code and generate executable file (the testcase exe)
    > gcc -o allocate_module_test allocate_module_test.c allocate_module.c -lcmocka -DUNIT_TESTING=1

    # run the test cases
    > allocate_module_test.exe
    [==========] tests: Running 2 test(s).
    [ RUN      ] leak_memory_test
    [  ERROR   ] --- Blocks allocated...
    allocate_module.c:42: note: block 000001ECB2D637B0 allocated here
    ERROR: leak_memory_test leaked 1 block(s)      <---------- report a mem leak

    [  FAILED  ] leak_memory_test
    [ RUN      ] no_leak_memory_test
    [       OK ] no_leak_memory_test
    [==========] tests: 2 test(s) run.
    [  PASSED  ] 1 test(s).
    [  FAILED  ] tests: 1 test(s), listed below:
    [  FAILED  ] leak_memory_test

    1 FAILED TEST(S)

It runs two cases and reports one issue of mem leak.

[cmocka-1.1.7 examples]: https://git.cryptomilk.org/projects/cmocka.git/tree/example?h=cmocka-1.1.7
