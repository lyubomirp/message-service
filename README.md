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
   "format":"text/html"
}
```
Where 
   - ```subject``` is the header of our message. REQUIRED FOR ALL MESSAGES
   - ```content``` is the content of our message. REQUIRED FOR ALL MESSAGES
   - ```type``` denotes the type of message being sent (i.e. slack, email, etc.) REQUIRED FOR ALL MESSAGES
   - ```recipients``` is an array of objects containing the name and contact of our... recipients (not required for slack messsages)
   - ```format``` is important mainly for emails, it's the format of the message (i.e. text/html, plain/text, multipart/alternative, etc.)

The recipients field can be omitted when sending a Slack message. We only need subject and content.
Any failed message will either be requeued, or logged in the DB if malformed and Nack'd from the queue.

Healthchecks for postgres and rabbitMQ are added on
```
/status
```
