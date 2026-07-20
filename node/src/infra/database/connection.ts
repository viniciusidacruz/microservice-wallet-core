import mysql from 'mysql2/promise'
import { env } from '../../env.js'

export async function createMySQLPool() {
  const pool = mysql.createPool({
    host: env.DB_HOST,
    port: env.DB_PORT,
    user: env.DB_USER,
    password: env.DB_PASSWORD,
    database: env.DB_NAME,
    waitForConnections: true,
    connectionLimit: 10,
  })

  for (let attempt = 1; attempt <= 30; attempt++) {
    try {
      const connection = await pool.getConnection()
      connection.release()
      console.log('connected to mysql')
      return pool
    } catch (error) {
      if (attempt === 30) {
        throw error
      }

      console.log(`waiting for mysql... attempt ${attempt}/30`)
      await new Promise((resolve) => setTimeout(resolve, 2000))
    }
  }

  return pool
}
