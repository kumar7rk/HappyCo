#conversation.go
check if a message is from a happyco user
yes
check their business id, message body
if the business id is any of tier 2 users
 call buildium.go
if message body contains password report or some certain phrases
 call password.go
and so on with other things you might want to handle in future
just add another if and the corresponding go file
and you should be good.

so yeah that' cool

now let's split the code into main.go and conversation.go

what does main.go do again please?
good question.
Just do the backend stuff
like handling webhooks routing connecting to db
maybe intercom connections?
