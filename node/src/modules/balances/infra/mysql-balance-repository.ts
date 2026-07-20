import type { Pool } from 'mysql2/promise'
import { Balance } from '../domain/balance.js'
import type { BalanceRepository } from '../domain/balance-repository.js'

type BalanceRow = {
  account_id: string
  balance: number
  updated_at: Date
}

export class MySQLBalanceRepository implements BalanceRepository {
  constructor(private readonly pool: Pool) {}

  async findByAccountId(accountId: string): Promise<Balance | null> {
    const [rows] = await this.pool.query(
      'SELECT account_id, balance, updated_at FROM balances WHERE account_id = ? LIMIT 1',
      [accountId],
    )

    const result = rows as BalanceRow[]
    const row = result[0]

    if (!row) {
      return null
    }

    return new Balance(row.account_id, Number(row.balance), new Date(row.updated_at))
  }

  async upsert(balance: Balance): Promise<void> {
    await this.pool.query(
      `
        INSERT INTO balances (account_id, balance, updated_at)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE
          balance = VALUES(balance),
          updated_at = VALUES(updated_at)
      `,
      [balance.accountId, balance.balance, balance.updatedAt],
    )
  }
}
