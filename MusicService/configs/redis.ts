import Redis from "ioredis";
import dotenv from "dotenv";

dotenv.config();

export async function Initialize() {
  // Buat instance Redis dengan konfigurasi
  const client = new Redis({
    host: process.env.REDIS_HOST || "localhost", // Default localhost jika REDIS_HOST tidak disediakan
    port: parseInt(process.env.REDIS_PORT || "6379", 10), // Default port 6379
    password: process.env.REDIS_PASSWORD || "", // Default tanpa password
    db: parseInt(process.env.REDIS_DB || "0", 10), // Default database index 0
    retryStrategy: (times) => {
      const delay = Math.min(times * 50, 2000); // Exponential backoff
      console.warn(`Retrying Redis connection in ${delay}ms...`);
      return delay;
    },
  });

  try {
    // Tes koneksi Redis
    await client.ping();
    console.log("Redis connected");
  } catch (error) {
    console.error("Redis connection error:", error);
  }

  return client;
}
