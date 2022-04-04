# Errandboi

Errandboi is a project that publishes given events to Nats and EMQX. It's basically a scheduler that releases the events based on their publish time.

## How to run your boy

This project requires MongoDB, Redis, EMQX and Nats JetStream.
You can run it via docker-compose and the following commands.
`make run all`
