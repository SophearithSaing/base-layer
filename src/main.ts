import { env } from './config/env.ts';
import { router } from './routes.ts';

Deno.serve({ port: env.PORT }, router);

console.log(`BaseLayer listening on http://localhost:${env.PORT}`);
