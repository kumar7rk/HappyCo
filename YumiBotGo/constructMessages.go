package main

func buildiumMessage(name string) (message string) {
	message := "Hi " + name + " 👋 \n \n <b>Our friends at Buildium support your Happy Inspector subscription. \n\n They can be reached at 888-414-1988, or by submitting a ticket through your Buildium account.</b>\n\nBuildium Support team are the best place to help you with this query as they understand your unique workflow and are trained in Happy Inspector 💫 \n \n Please also feel free to take a look through our FAQ on the Buildium integration:  \n https://intercom.help/happyco/frequently-asked-questions/buildium-integration-faq/faq-buildium-integration  \n Thanks!  \n HappyBot ☺"
}

func passwordMessage(name string) (message string) {
	message := "Hi " + name + " 👋 \n \n It looks like you might be having trouble logging in? \n\n You can reset your password by entering your email <a href='https://manage.happyco.com/password/forgot'> here </a> \n \n Thanks!  \n HappyBot ☺ \n\n <i>Need to contact a human....... just reply</i>"
}

func followUpMessage(name string) (message string) {
	message := "Hi " + name + " 👋 \n \n Just wanted to message and check in, and see how you're going? \n\n Did our last message help? \n \n Shout out if you need any help.\n"+<author's name>  ☺ "
}

func closingMessage(name string) (message string) {

}