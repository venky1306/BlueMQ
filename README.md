## Overview


BlueMQ is in-memory publish-subscribe messaging system written in Go, which can provides instant event notification for distributed systems. Supports push-based systems where publisher distributes messages to all subscribers when an event occurs.

BlueMQ is designed as an in-memory database for low-latency data transfer between applications.

## Architecture

Delivers the message to all the connected subscribers by checking the message topic. Consumers connect to the relevant topic and extract data from its partition.


## Benchmark
Single producer single consumer BlueMQ performed with Ultra-low latency of 1ms when distributing 100KB sized messages.Higher throughput of 380MB/s.