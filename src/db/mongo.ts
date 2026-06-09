import { MongoClient } from 'mongodb';
import { env } from '../config/env.ts';

const client = new MongoClient(env.MONGODB_URI);
export const db = client.db(env.MONGODB_DB_NAME);

export async function connectMongo() {
  await client.connect();
  await db.command({ ping: 1 });

  console.log(`Mongo connected: ${env.MONGODB_DB_NAME}`);
}
