import type { FastifyInstance } from 'fastify'
import type { GetBalanceUseCase } from '../application/get-balance-use-case.js'
import { AppError } from '../../../shared/errors/app-error.js'

export function balanceRoutes(
  app: FastifyInstance,
  getBalanceUseCase: GetBalanceUseCase,
): void {
  app.get<{ Params: { account_id: string } }>(
    '/balances/:account_id',
    async (request, reply) => {
      try {
        const output = await getBalanceUseCase.execute(request.params.account_id)
        return reply.status(200).send(output)
      } catch (error) {
        if (error instanceof AppError) {
          return reply.status(error.statusCode).send({ error: error.message })
        }

        request.log.error(error)
        return reply.status(500).send({ error: 'internal server error' })
      }
    },
  )
}
