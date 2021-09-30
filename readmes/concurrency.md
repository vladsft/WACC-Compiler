# Concurrency

We have chosen to base our concurrency model off Go. Go has goroutines, which are started with the following syntax `go function(args)`. This starts a green thread which runs the function with the given arguments. 

This syntax doesn't allow use to directly return the results of the function though. Go solved this problem by introducing channels, which can be used to send the result of a goroutine. Channels are thread safe and thus a better option that sending pointers to the function.

Our implementation will be based off this model, but due to time constraints we won't be able to introduce green threads, and instead will be using normal threads.

## Syntax

We introduced the `wacc` keyword for running a function in a separate thread.
`wacc function(args)`

We also added the `lock` type which is a simple mutex.
`lock l` - declares a lock
`acquire l` - acquires the lock l
`release l` - releases the lock l
`tryLock l` - locks the lock if it can
`free l` - frees the lock l

We will add channels via a standard library rather than implementing them directly in assembly.

## Semantics

`wacc` is a statement and so cannot be on the right hand side of an assignment. Apart from that it has similar semantics to `call`.

Locks can only be used with the keywords `acquire`, `release` and `free`.

## Code Generation

We will be using the pthreads library to start new threads and implement lock methods.

`wacc function(args)` roughly translates to this C code:
```
pthread_t *t = malloc(sizeof(pthread_t));
void* argMem = malloc(sizeof(args));
pthread_create(t, NULL, function, argMem);
pthread_detach(thread_id);
free(t);
```
The function `function` will also have a header for concurrent calls which would put the arguments on the stack, and free `argMem`.

`lock l` roughly translates to this C code:
`pthread_mutex_t l;`

`acquire l` roughly translates to this C code:
`pthread_mutex_lock(&l)`

`release l` roughly translates to this C code:
`pthread_mutex_unlock(&l)`

`tryLock l` roughly translates to this C code:
`pthread_mutex_trylock(&l)`

`free l` roughly translates to this C code:
`pthread_mutex_destroy(&l)`
