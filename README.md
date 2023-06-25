# Message Service

A simple messaging service that connects to a RabbitMQ queue and consumes messages, and sends them via Slack/Email.

Expected message format is

```
{
   "recipients":[
      {
         "name":"Test",
         "contact":"someemail@abv.bg"
      }
   ],
   "subject":"Testing",
   "content":"Test Message Body",
   "type":"email",
   "format":"plain/text"
}
```
The recipients field can be omitted when sending a Slack message. We only need subject and content.
Any failed message with either be requeued, or logged in the DB if malformed, and Nack'd from the queue.

Healthchecks for postgres and rabbitMQ are added on
```
/status
```
