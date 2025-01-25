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

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR]consumer: %s", err.Error())

			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
		// добавьте условие выхода из цикла
	}
}
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("получено новое событие '%s'", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("низя обработать: '%s' ", err.Error())
			continue
		}
	}
	return nil
}
