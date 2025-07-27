Rest API -> Kafka
multiple consumer to consume kafka and publish to operators
Express SMS

-- -
OS
virtualization -> resource -> response time == higher priority
customers with high rate of SMS go to partitions with high amount of wait
100,000
10

-- -
Guarantee 5 min to receive the sms to user for express service
Dynamicly allocate pods by topic size
Monitor the amount of time a consumer process a message, then use this average amount to calculate amount time is needed to process all messages.
avg time to process all messages = (Avg Processing Time) * (count of messages) / (count of consumers)

-- -
How to handle budgets
Massive amount of update on database if we want to update all of the events

Use a RDBMS to store values and a cache like redis with write through cache architecture.

-- -
Reports
Timeseries
Clickhouse

-- -
Retry policy
on the kafka side if the consumer does not ack on timeout, then retry message to another consumer
On the consumer side if the provider is not available retry by failing the message and requeueing the message on the topic (if it's possible to be sent on the back on the topic is better)

-- -


Kubernetes -> Auto scale
Saga pattern budget management

-- -
100,000,000 SMS per day
100,000,000 / (24 * 60 * 60) * 1.3 = 1,500 SMS Per Second
Bulk request for SMSs
Single request but make 

-- -
Stashing in each 10 seconds

- For budget managment
	- Stashing -> all the requests from user are accepted  and stashed for each `user_id` and `formating` and for each n seconds are sent as a bulk request to budget management.
		- Durability
		- Store on Redis
	- Cache budget management data in API Gateway
		- It's very error prone
- For Priority
	- Normal
		- Fair distribution of resource to customers
		- Multiple queue
			- High
			- Normal
	- Express
		- SLA

Redis stream
DragonFly
Kafka 
When there is multiple 

Is it possible to implement consumer orchestrator with K8S or goroutines


