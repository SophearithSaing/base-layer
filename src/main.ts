import { router } from './routes.ts';

const port = Number(Deno.env.get('PORT') ?? 8000);

Deno.serve({ port }, router);

console.log(`BaseLayer listening on http://localhost:${port}`);
