package queue

import "encoding/json"

func (q *QueueManager) Publish(evt QueueEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	return q.client.RPush(q.ctx, q.queue, data).Err()
}
