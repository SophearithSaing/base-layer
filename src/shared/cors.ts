import { env } from '../config/env.ts';

const allowedOrigin = [env.CLIENT_ORIGIN];

export function getCorsHeaders(req: Request): Headers {
  const headers = new Headers();
  const origin = req.headers.get('origin');

  if (origin && allowedOrigin.includes(origin)) {
    headers.set('Access-Control-Allow-Origin', origin);
    headers.set('Access-Control-Allow-Credentials', 'true');
  }

  headers.set('Access-Control-Allow-Methods', 'GET,POST,PATCH,DELETE,OPTIONS');
  headers.set('Access-Control-Allow-Headers', 'content-type');

  return headers;
}

export function preflight(req: Request): Response {
  return new Response(null, {
    status: 204,
    headers: getCorsHeaders(req),
  });
}
