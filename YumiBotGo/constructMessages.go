package main

func buildiumMessage(name string) (message string) {
	message = "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"
	return
}

func passwordMessage(name string) (message string) {
	message = "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to contact a human....... just reply</i>"
	return
}

func followUpMessage(name string, authorName string) (message string) {
	message = "Hey " + name + " 👋 \n \n Just wanted to message and check in, and see how you're going? \n\n Did our last message help? \n \n Shout out if you need any help.\n\n" + authorName + " 😄 "
	return
}

func closingMessage(name string, authorName string) (message string) {
	message = "Hey " + name + " 👋 \n \n We're still here if you need any help. \n\n I'm closing this conversation for now, but you can reopen it at anytime by replying. \n \n We value your feedback Rohit- please rate us on AppStore: https://hpy.io/appstore-review 💯  \n\n" + authorName + " 😄 \n\n Get real-time intelligence on property conditions and portfolio trends to optimize operations, achieve higher NOI, and make better business decisions: http://hpy.io/happy-insights"
	return
}

func welcomeMessage(name string) (message string) {
	message = "Hey " + name + " 👋 \n \n I'm your friendly HappyBot 😎 \n\n We've received your message and the team will respond as soon as possible 🕛 🔜  \n \n In the meantime, check out our iOS user manual hpy.io/manual 📓  or search the issue on http://support.happy.co. 🌏  \n\n Trouble seeing your reports/data? Make sure your Sync/Cloud tab in app is clear"
	return
}
