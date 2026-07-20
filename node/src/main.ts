import { env } from './env.js'
import { createMySQLPool } from './infra/database/connection.js'
import { createHttpServer } from './infra/http/server.js'
import { createKafka } from './infra/kafka/kafka.js'
import { GetBalanceUseCase } from './modules/balances/application/get-balance-use-case.js'
import { UpsertBalancesFromEventUseCase } from './modules/balances/application/upsert-balances-from-event-use-case.js'
import { BalanceUpdatedKafkaConsumer } from './modules/balances/infra/balance-updated-consumer.js'
import { MySQLBalanceRepository } from './modules/balances/infra/mysql-balance-repository.js'

async function bootstrap() {
  const pool = await createMySQLPool()
  const balanceRepository = new MySQLBalanceRepository(pool)

  const getBalanceUseCase = new GetBalanceUseCase(balanceRepository)
  const upsertBalancesFromEventUseCase = new UpsertBalancesFromEventUseCase(
    balanceRepository,
  )

  await createHttpServer(getBalanceUseCase)

  const kafka = createKafka()
  const consumer = kafka.consumer({ groupId: env.KAFKA_GROUP_ID })

  await consumer.connect()
  console.log('connected to kafka')

  const balanceUpdatedConsumer = new BalanceUpdatedKafkaConsumer(
    consumer,
    env.KAFKA_TOPIC_BALANCE_UPDATED,
    upsertBalancesFromEventUseCase,
  )

  await balanceUpdatedConsumer.start()
  console.log(`consuming topic: ${env.KAFKA_TOPIC_BALANCE_UPDATED}`)
}

bootstrap().catch((error) => {
  console.error('failed to start balances service:', error)
  process.exit(1)
})
