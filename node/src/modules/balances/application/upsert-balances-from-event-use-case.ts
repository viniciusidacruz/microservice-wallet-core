import { Balance } from '../domain/balance.js'
import type { BalanceRepository } from '../domain/balance-repository.js'

export interface BalanceUpdatedPayload {
  account_from_id: string
  account_to_id: string
  balance_account_from: number
  balance_account_to: number
}

export class UpsertBalancesFromEventUseCase {
  constructor(private readonly balanceRepository: BalanceRepository) {}

  async execute(payload: BalanceUpdatedPayload): Promise<void> {
    const accountFrom = Balance.create(
      payload.account_from_id,
      payload.balance_account_from,
    )
    const accountTo = Balance.create(
      payload.account_to_id,
      payload.balance_account_to,
    )

    await this.balanceRepository.upsert(accountFrom)
    await this.balanceRepository.upsert(accountTo)
  }
}
