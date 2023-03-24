# go-streaming-with-kafka
Go agent that report to a centralized data platform via streaming


Having a Go agent allows us to utilize the same code base for various operating systems, and the fact that it has good 
`integration with JSON data` and packages such as a `MongoDB driver` and `Confluent Go Kafka Client` makes it a compelling candidate for the presented use case.

Demonstrates how file size data on a host is monitored from a cross-platform agent written in Golang via a Kafka cluster using a Confluent hosted sink connector to MongoDB Atlas. MongoDB Atlas stores the data in a time series collection. The MongoDB Charts product is a convenient way to show the gathered data to the user.
