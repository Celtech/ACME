package queue

func (q *QueueManager) Publish(message string) error {
	return q.client.RPush(q.ctx, q.queue, message).Err()
}
