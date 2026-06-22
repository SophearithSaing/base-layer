import { aiRoutes } from './ai/ai.routes.ts';
import { authRoutes } from './auth/auth.routes.ts';
import { projectRoutes } from './projects/projects.routes.ts';
import { addCorsHeaders, preflight } from './shared/cors.ts';
import { error, Handler, HttpError, HttpMethod } from './shared/http.ts';

export type Route = {
  method: HttpMethod;
  pattern: URLPattern;
  handler: Handler;
};

const routes: Route[] = [
  ...authRoutes,
  ...projectRoutes,
  ...aiRoutes,
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/' }),
    handler: () =>
      Response.json({
        name: 'BaseLayer',
        status: 'ok',
      }),
  },
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/health' }),
    handler: () =>
      Response.json({
        ok: true,
        timestamp: new Date().toISOString(),
      }),
  },
];

export async function router(req: Request): Promise<Response> {
  if (req.method === 'OPTIONS') {
    return preflight(req);
  }

  const url = new URL(req.url);

  for (const route of routes) {
    if (route.method !== req.method) continue;

    const match = route.pattern.exec({
      pathname: url.pathname,
    });

    if (!match) continue;

    const groups = (match.pathname.groups ?? {}) as Record<string, string>;

    try {
      const res = await route.handler(req, groups);
      return addCorsHeaders(req, res);
    } catch (err) {
      if (err instanceof HttpError) {
        return addCorsHeaders(req, error(err.message, err.status));
      }

      console.error(err);
      return addCorsHeaders(req, error('Internal server error', 500));
    }
  }

  return addCorsHeaders(req, error('Not Found', 404));
}
