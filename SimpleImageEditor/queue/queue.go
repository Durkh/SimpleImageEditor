package queue

type Queue chan rune

func New(bufferSize uint) (Queue, func()) {

	q := make(chan rune, bufferSize)

	return q, func() {
		close(q)
	}
}

func (q Queue) Sender() chan<- rune {

	return q
}

func (q Queue) Receiver() <-chan rune {

	return q
}
