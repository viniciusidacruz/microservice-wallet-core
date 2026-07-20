import type { Consumer } from 'kafkajs'
import type { UpsertBalancesFromEventUseCase } from '../application/upsert-balances-from-event-use-case.js'

type BalanceUpdatedEvent = {
  name: string
  payload: {
    account_from_id: string
    account_to_id: string
    balance_account_from: number
    balance_account_to: number
  }
}

export class BalanceUpdatedKafkaConsumer {
  constructor(
    private readonly consumer: Consumer,
    private readonly topic: string,
    private readonly upsertBalancesFromEventUseCase: UpsertBalancesFromEventUseCase,
  ) {}

  async start(): Promise<void> {
    await this.consumer.subscribe({ topic: this.topic, fromBeginning: true })

    await this.consumer.run({
      eachMessage: async ({ message }) => {
        if (!message.value) {
          return
        }

        try {
          const event = JSON.parse(message.value.toString()) as BalanceUpdatedEvent

          if (!event.payload) {
            console.warn('ignoring kafka message without payload')
            return
          }

          await this.upsertBalancesFromEventUseCase.execute(event.payload)

          console.log('balance updated from kafka event:', {
            account_from_id: event.payload.account_from_id,
            account_to_id: event.payload.account_to_id,
            balance_account_from: event.payload.balance_account_from,
            balance_account_to: event.payload.balance_account_to,
          })
        } catch (error) {
          console.error('failed to process balance_updated event:', error)
        }
      },
    })
  }
}
