import { env } from './config/env.ts';
import { ensureIndexes } from './db/collections.ts';
import { connectMongo } from './db/mongo.ts';
import { router } from './routes.ts';

await connectMongo();
await ensureIndexes();

Deno.serve({ port: env.PORT }, router);
