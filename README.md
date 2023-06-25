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

