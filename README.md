# go-cancelonerrorgroup

This is a little example package that I wrote because I misunderstood how the `errgroup.Group` works. 

## errgroup.Group

When I first encountered the `errgroup.Group` I thought it worked something like this:

1. create `errgroup.Group`
2. throw it some goroutines and wait
3. if any of them finish with an error anything still running **stops**
4. return and handle error

This was useful as I can kick off a bunch of goroutines independantly working away but then early exit if I can't complete one of the requests (I'm thinking something like calling a bunch of services to build a response). 

That's not, however, how `errgroup.Group` works. Instead an `errgroup.Group` behaves similarly to the `sync.WaitGroup` which means that it wait for all running gorountines to finish before moving forward (this is probably the right behaviour, in most cases, as it would allow in-flight work to complete).  

## CancelOnErrorGroup

So what if I just wrote the behaviour I do want?

That's what `CancelOnErrorGroup` aims to be. In short it's a wrapper around the `sync.WaitGroup` that allows for an underlying context to be cancelled if any of the gorountines in the group error.

It works as described above and first collects a number of `handles` which can then passed to the `sync.WaitGroup` when the `ceg.Wait(ctx)` is called. 

Each of the `handles` accepts a `context` and this is used to drive "cancel on error" behaviour as when a goroutine errors the error is stored, the context is cancelled and the `wg.Wait()` will end. Returning the error keeps the implementation similar to that of the `errgroup.Group`. 


