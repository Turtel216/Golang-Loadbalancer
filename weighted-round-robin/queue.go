package main

// add server to queue
func enqueue(queue []Server, element Server) []Server {
	queue = append(queue, element)

	return queue
}

// pop server from queue
func dequeue(queue []Server) (Server, []Server) {
	element := queue[0]

	// The first element is the one to be dequeue
	if len(queue) == 1 {
		var tmp = []Server{}
		return element, tmp
	}

	return element, queue[1:]
}
