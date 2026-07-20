import { Kafka } from 'kafkajs'
import { env } from '../../env.js'

export function createKafka() {
  return new Kafka({
    clientId: env.KAFKA_CLIENT_ID,
    brokers: env.KAFKA_BROKERS.split(',').map((broker) => broker.trim()),
    retry: {
      retries: 30,
      initialRetryTime: 1000,
    },
  })
}
