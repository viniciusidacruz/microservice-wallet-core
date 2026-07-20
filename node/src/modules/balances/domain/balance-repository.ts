import type { Balance } from './balance.js'

export interface BalanceRepository {
  findByAccountId(accountId: string): Promise<Balance | null>
  upsert(balance: Balance): Promise<void>
}
