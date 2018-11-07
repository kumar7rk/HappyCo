package main

func buildiumMessage(name string) (message string) {
	message = "Hi " + name + " 👋 \n \n Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at <b>888-414-1988</b>, or by <b>submitting a ticket</b> through your Buildium account.\n\nThey understand your workflow and are trained in Happy Inspector. \n \n Check out our <a href=\"https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration\">FAQ on the Buildium integration</a>\n \n HappyBot 😎"
	return
}

func passwordMessage(name string) (message string) {
	message = "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot 😎 \n\n <i>Need to speak to a human....... just reply</i>"
	return
}

func followUpMessage(name string, authorName string) (message string) {
	message = "Hey " + name + " 👋 \n \n Just wanted to message and check in, and see how you're going? \n\n Did our last message help? \n \n Shout out if you need any help.\n\n" + authorName + " 😄 "
	return
}

func closingMessage(name string, authorName string) (message string) {
	message = "Hey " + name + " 👋 \n \n We're still here if you need any help. \n\n I'm closing this conversation for now, but you can reopen it at anytime by replying. \n \n We value your feedback " + name + " - please rate us on AppStore: https://hpy.io/appstore-review or Google PlayStore: https://hpy.io/get-android 💯 \n\n" + authorName + " 😄 \n\n Get real-time intelligence on property conditions and portfolio trends to optimize operations, achieve higher NOI, and make better business decisions: https://hpy.io/happy-insights"
	return
}

func welcomeMessage(name string) (message string) {
	message = "Hey " + name + " 👋 \n \n I'm your friendly HappyBot 😎 \n\n We've received your message and the team will respond as soon as possible 🕛 🔜  \n \n <b>Need help faster?</b> check out our iOS user manual https://hpy.io/manual 📓  or search the issue on https://support.happy.co 🌏 \n\n <i>Trouble seeing your reports/data? Make sure your Sync/Cloud tab in app is clear  or try these</i> <a href=\"https://support.happy.co/further-help-and-troubleshooting\">troubeshooting tips</a> 👍"
	return
}

func signUpMessage(name string) (message string) {
	message = "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n What email are you using to login? \n\n If you're trying to start a free trial, please register here: <a href='https://manage.happyco.com/signup'> https://manage.happyco.com/signup </a> \n \n Thanks! \n HappyBot 😎 \n\n <i>Need to speak to a human....... just reply</i>"
	return
}
