#main.go
just to handle the db and route

#new conversation handler 
- this will be called when the corresponding webhook runs (?)

- parse the json
- just check if the message is from a user or lead
- if lead forget about it
- if user call new conversation trigger

new conversation trigger
- processing json
- kind of like current makeAndSendNote func
- so basically it's a bunch of ifs - if the user is buildium, if the message is password error report
  - decide the sequence/priority

If the condition allow (message from a buildium user for example)
run a command.

#Command

command is a file
an example is a Buildium.go
whenever it's called, its job is to send a message to the user
both the new conversation and the yumi run buildium command will call this file

This is the first module




The second module is note handler

we are having different webhooks for new conversation and note handler and for any other future feature

#Note handler 
is a different let's call it module
bascially has two different parts (not sure what's the right word here)
run and help

this bit honestly looks simple 
since we're having a seperate webhook or even if have to check if this is a note
we can just check what's after yumi run and call that file name for example.
or maybe a bunch of ifs

###
what Michael wants me to do by Tuesday
Put the main.go in a separate file
have a new file named Conversation
have a function parsing json if user call the next one (like makeAndSendNote) if not don't do anything

in this function check if the message is from a Buildium user and maybe also read the password error report
make new files- buildium, password reset
call them Boom