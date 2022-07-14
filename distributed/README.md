# Distributed System

A distributed system contains services running
inside a network instead of in a standalone computer.

This project shows how to design the system structure
and how to write registry service to support services
(de)registration, notification and status monitoring.

Structure of the system:

```golang
                                                              -----------
                              ----------------------------->  |         |
                              |                               |   Log   |
                              |                      ======>  | Service |
                              |                      ||       -----------
                -----------   |     ------------     ||
                |         | ---     |          | <====
    User  <---> | Web UI  | <=====> | Registry |
    App         | Service | ---     | Service  | <====
                -----------   |     ------------     ||
                              |                      ||       -----------
                              |                       =====>  |         |
                              |                               | Grading |
                              ----------------------------->  | Service |
                                                              -----------
```

Structure of registry service:

```golang
    --------------                          ------------------
    |            |  <--- (de)register ----  |                | ------
    |            |  ----    notify    --->  |    service 1   |      |
    |            |  ----   heartbeat  --->  | (require 2, x) | ---  |
    |            |                          ------------------   |  |
    |            |                                               |  |
    |            |                          ------------------   |  |
    |            |  <--- (de)register ----  |                | <--  |
    |            |  ----    notify    --->  |    service 2   |      |
    |            |  ----   heartbeat  --->  |   (require y)  |      |
    |            |                          ------------------      |
    |  Registry  |                                   .              |
    |            |                                   .              |
    |   Service  |                                   .              |
    |            |                          ------------------      |
    |            |  <--- (de)register ----  |                | <-----
    |            |  ----    notify    --->  |    service x   |
    |            |  ----   heartbeat  --->  |    (require N) | ----
    |            |                          ------------------    |
    |            |                                                |
    |            |                          ------------------    |
    |            |  <--- (de)register ----  |                |    |
    |            |  ----              ----  |    service N   | <---
    |            |  ----   heartbeat  --->  |                |
    --------------                          ------------------

```