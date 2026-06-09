import { env } from './config/env.ts';
import { connectMongo } from './db/mongo.ts';
import { router } from './routes.ts';

await connectMongo();

Deno.serve({ port: env.PORT }, router);
