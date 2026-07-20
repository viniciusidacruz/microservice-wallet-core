export class Balance {
  constructor(
    public readonly accountId: string,
    public balance: number,
    public updatedAt: Date,
  ) {}

  static create(accountId: string, balance: number): Balance {
    return new Balance(accountId, balance, new Date())
  }

  update(balance: number): void {
    this.balance = balance
    this.updatedAt = new Date()
  }
}
