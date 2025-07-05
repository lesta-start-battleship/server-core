package event

type MatchEventDispatcher struct {
	publisher MatchEventPublisher
}

func NewMatchEventDispatcher(publisher MatchEventPublisher) *MatchEventDispatcher {
	return &MatchEventDispatcher{publisher: publisher}
}

func (d *MatchEventDispatcher) DispatchMatchResult(result MatchResult) error {
	return d.publisher.PublishMatchResult(result)
}

func (d *MatchEventDispatcher) DispatchUsedItem(item Item) error {
	return d.publisher.PublishUsedItem(item)
}
