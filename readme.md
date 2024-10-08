Initially, Redis In GO, hence name: rigo.
Now, an all new queueing system

To test current queue system, run:
`go run example/main.go`

Ideas:

- Have a setup where the queue can be hosted as a separate standalone server operated via a tcp client.
- Have a setup where the queue can also be installed into an application and be used locally, managed completely by the devs
- Implement message limit per channel (buffering)
- Record connection count
