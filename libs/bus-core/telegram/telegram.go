package telegram

const GenerateSubscriptionUrlSubject = "telegram.generate_subscription_url"

type GenerateSubscriptionUrlInput struct {
	Title      string
	StarsPrice int
}
