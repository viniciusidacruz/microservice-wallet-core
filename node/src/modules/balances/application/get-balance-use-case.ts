import { NotFoundError } from '../../../shared/errors/app-error.js'
import type { BalanceRepository } from '../domain/balance-repository.js'

export interface GetBalanceOutput {
  account_id: string
  balance: number
}

export class GetBalanceUseCase {
  constructor(private readonly balanceRepository: BalanceRepository) {}

  async execute(accountId: string): Promise<GetBalanceOutput> {
    const balance = await this.balanceRepository.findByAccountId(accountId)

    if (!balance) {
      throw new NotFoundError('balance not found')
    }

    return {
      account_id: balance.accountId,
      balance: balance.balance,
    }
  }
}
