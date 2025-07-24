package activity

import (
	"context"
	"fmt"
)

type GreetingActivities struct {
	LocalGreeting string
}

func (s *GreetingActivities) Greet(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("%s, %s!", s.LocalGreeting, name)
	return greeting, nil
}

func (s *GreetingActivities) SayGoodbye(ctx context.Context) (string, error) {
	return "Goodbye!", nil
}
