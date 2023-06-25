# Message Service

A simple messaging service that connects to a RabbitMQ queue and consumes messages, and sends them via Slack/Email.

Expected message format is

```
{
   "recipients":[
      {
         "name":"Test",
         "contact":"lpetkov44@yahoo.com"
      }
   ],
   "subject":"Testing",
   "content":"Test Message Body",
   "type":"slack"
}
```
Any failed message with either be requeued, or logged in the DB if malformed, and Nack'd from the queue.

Healthchecks for postgres and rabbitMQ are added on
```
/status
```
