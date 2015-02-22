package models

import "container/list"

const (
	ChannelCapacity = 10
	Success         = "SUCCESS"
	Failure         = "FAILURE"
)

var quizMonitor *QuizMonitor

type QuizMonitor struct {
	started bool

	subscribers *list.List
	archive     *list.List

	subscribeChannel   chan (chan Subscription)
	unsubscribeChannel chan Subscription
	publishChannel     chan AttemptWrapper

	controlChannel chan struct{}
}

func QuizMonitorInstance() *QuizMonitor {
	if quizMonitor == nil {
		q := &QuizMonitor{}

		q.subscribers = list.New()
		q.archive = list.New()

		q.subscribeChannel = make(chan (chan Subscription), ChannelCapacity)
		q.unsubscribeChannel = make(chan Subscription, ChannelCapacity)
		q.publishChannel = make(chan AttemptWrapper, ChannelCapacity)
		q.controlChannel = make(chan struct{}, 1)

		quizMonitor = q
	}
	return quizMonitor
}

func (q *QuizMonitor) Start() {
	if !q.started {
		go q.start()
	}
}

func (q *QuizMonitor) Stop() {
	if q.started {
		q.controlChannel <- struct{}{}
	}
}

func (q *QuizMonitor) Subscribe() Subscription {
	subscription := make(chan Subscription)
	q.subscribeChannel <- subscription
	return <-subscription
}

func (q *QuizMonitor) Publish(status string, attempt Attempt) {
	var wrapper AttemptWrapper
	switch status {
	case Success:
		wrapper = AttemptWrapper{Success: true, Attempt: attempt}
	case Failure:
		wrapper = AttemptWrapper{Success: false, Attempt: attempt}
	}
	q.publishChannel <- wrapper
}

func (q *QuizMonitor) Unsubscribe(subscription Subscription) {
	q.unsubscribeChannel <- subscription
	q.drain(subscription.New)
}

func (q *QuizMonitor) start() {
	defer close(q.unsubscribeChannel)
	defer close(q.subscribeChannel)
	defer close(q.publishChannel)
	q.started = true
	for {
		select {
		case s := <-q.subscribeChannel:
			subscriber := make(chan AttemptWrapper, ChannelCapacity)
			for a := q.archive.Front(); a != nil; a = a.Next() {
				subscriber <- a.Value.(AttemptWrapper)
			}
			q.subscribers.PushBack(subscriber)
			s <- Subscription{New: subscriber}
		case r := <-q.publishChannel:
			for s := q.subscribers.Front(); s != nil; s = s.Next() {
				s.Value.(chan AttemptWrapper) <- r
			}
			if q.archive.Len() > ChannelCapacity {
				q.archive.Remove(q.archive.Front())
			}
			q.archive.PushBack(r)
		case u := <-q.unsubscribeChannel:
			for s := q.subscribers.Front(); s != nil; s = s.Next() {
				if s.Value.(chan AttemptWrapper) == u.New {
					q.subscribers.Remove(s)
					break
				}
			}
		case <-q.controlChannel:
			return
		}
	}
}

func (q QuizMonitor) drain(ch chan AttemptWrapper) {
	for {
		defer close(ch)
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
