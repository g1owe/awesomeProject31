package event_consumer

import (
	"awesomeProject3/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func (c *Consumer) Start() error {
	var stop bool
	for !stop {
		getEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR]consumer: '%s' ", err)
			continue
		}
		if len(getEvents) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}
		if err := c.handleEvent(getEvents); err != nil {
			log.Println(err)
			continue
		}
		// добавьте условие выхода из цикла
		if len(getEvents) == 0 {
			stop = true
		}
	}
	return nil
}
func (c *Consumer) handleEvent(events []events.Event) error {
	for _, event := range events {
		log.Printf("получено новое событие '%s'", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("низя обработать: '%s' ", err)
			continue
		}
	}
	return nil
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}
