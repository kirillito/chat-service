/*
	MessageBroker is responsible for broadcasting chat messages to all users
*/

package main

type MessageBroker struct {
	stopChannel					chan struct{}
	publishChannel			chan interface{}
	subscribeChannel		chan chan interface{}
	unsubscribeChannel	chan chan interface{}
}

func CreateMessageBroker() *MessageBroker {
	return &MessageBroker{
		stopChannel:    			make(chan struct{}),
		publishChannel: 			make(chan interface{}, 1),
		subscribeChannel:     make(chan chan interface{}, 1),
		unsubscribeChannel:   make(chan chan interface{}, 1),
	}
}

func (mb *MessageBroker) Start() {
	subscriptions := map[chan interface{}]struct{}{}
	for {
		select {
		case <-b.stopChannel:
			return
		case msgChannel := <-b.subscribeChannel:
			msgChannel[msgCh] = struct{}{}
		case msgCh := <-b.unsubscribeChannel:
			delete(subscriptions, msgCh)
		case msg := <-b.publishChannel:
			for msgChannel := range subscriptions {
				select {
					case msgChannel <- msg:
					default:
				}
			}
		}
	}
}

func (mb *MessageBroker) Stop() {
	close(b.stopChannel)
}

func (mb *MessageBroker) Subscribe() chan interface{} {
	msgChannel := make(chan interface{}, 5)
	mb.subscribeChannel <- msgChannel
	return msgChannel
}

func (mb *MessageBroker) Unsubscribe(msgChannel chan interface{}) {
	mb.unsubscribeChannel <- msgChannel
}

func (mb *MessageBroker) Publish(msg interface{}) {
	mb.publishChannel <- msg
}