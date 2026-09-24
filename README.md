A simple system that helps you AI-review your PRs :)


Alright how to use it?


This is our env file:

// -- This is the url of your git provider. For example, https://github.com, or https://felan.gitea.ir

GIT_URL=

// -- Your Git access token. You can obtain it by the applications section in your git profile's privacy options

GIT_TOKEN=

// -- This is the url of your OLLAMA instance. For example: http://127.0.0.1:11434

MODEL_URL=


Fill the env file and then run the system using:

go run main.go

And how to run it in production?

Put it on your go-supporting server, and then:

source start.sh --profile prod --network full --action up


Your system is operational.

Any contribution is welcomed!
