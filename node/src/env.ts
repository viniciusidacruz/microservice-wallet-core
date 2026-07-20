import { z } from 'zod'

const envSchema = z.object({
  NODE_ENV: z.enum(['development', 'production', 'test']).default('development'),
  HTTP_PORT: z.coerce.number().default(3003),
  DB_HOST: z.string().default('localhost'),
  DB_PORT: z.coerce.number().default(3307),
  DB_USER: z.string().default('root'),
  DB_PASSWORD: z.string().default('root'),
  DB_NAME: z.string().default('balances'),
  KAFKA_BROKERS: z.string().default('localhost:9092'),
  KAFKA_CLIENT_ID: z.string().default('balances-service'),
  KAFKA_GROUP_ID: z.string().default('balances-consumer-group'),
  KAFKA_TOPIC_BALANCE_UPDATED: z.string().default('balance_updated'),
})

export const env = envSchema.parse(process.env)
