import Fastify from 'fastify'
import type { GetBalanceUseCase } from '../../modules/balances/application/get-balance-use-case.js'
import { balanceRoutes } from '../../modules/balances/infra/balance-routes.js'
import { env } from '../../env.js'

export async function createHttpServer(getBalanceUseCase: GetBalanceUseCase) {
  const app = Fastify({
    logger: true,
  })

  app.get('/health', async () => ({ status: 'ok' }))

  balanceRoutes(app, getBalanceUseCase)

  await app.listen({
    port: env.HTTP_PORT,
    host: '0.0.0.0',
  })

  console.log(`balances service listening on port ${env.HTTP_PORT}`)

  return app
}
