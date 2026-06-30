import { env } from '../config/env.ts';

const allowedOrigins = env.CLIENT_ORIGINS.split(',')
  .map((origin) => origin.trim())
  .filter(Boolean);

export function getCorsHeaders(req: Request): Headers {
  const headers = new Headers();
  const origin = req.headers.get('origin');

  if (origin && allowedOrigins.includes(origin)) {
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

export function addCorsHeaders(req: Request, res: Response): Response {
  const headers = new Headers(res.headers);

  getCorsHeaders(req).forEach((value, key) => {
    headers.set(key, value);
  });

  return new Response(res.body, {
    status: res.status,
    statusText: res.statusText,
    headers,
  });
}
