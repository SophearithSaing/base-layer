import { assertEquals } from '@std/assert';

async function importRouter() {
  Deno.env.set('MONGODB_URI', 'mongodb://localhost:27017');
  Deno.env.set('JWT_SECRET', 'test-secret-at-least-32-characters');
  Deno.env.set('CLIENT_ORIGIN', 'http://localhost:5173');

  return await import('../src/routes.ts');
}

Deno.test('router adds CORS headers to successful responses', async () => {
  const { router } = await importRouter();
  const res = await router(
    new Request('http://localhost:8000/', {
      headers: { origin: 'http://localhost:5173' },
    }),
  );

  assertEquals(res.status, 200);
  assertEquals(
    res.headers.get('access-control-allow-origin'),
    'http://localhost:5173',
  );
});

Deno.test('router adds CORS headers to not found responses', async () => {
  const { router } = await importRouter();
  const res = await router(
    new Request('http://localhost:8000/missing', {
      headers: { origin: 'http://localhost:5173' },
    }),
  );

  assertEquals(res.status, 404);
  assertEquals(
    res.headers.get('access-control-allow-origin'),
    'http://localhost:5173',
  );
  assertEquals(res.headers.get('access-control-allow-credentials'), 'true');
});

Deno.test('router adds CORS headers to protected endpoint errors', async () => {
  const { router } = await importRouter();
  const res = await router(
    new Request('http://localhost:8000/projects', {
      headers: { origin: 'http://localhost:5173' },
    }),
  );

  assertEquals(res.status, 401);
  assertEquals(
    res.headers.get('access-control-allow-origin'),
    'http://localhost:5173',
  );
});

Deno.test('router handles CORS preflight', async () => {
  const { router } = await importRouter();
  const res = await router(
    new Request('http://localhost:8000/projects', {
      method: 'OPTIONS',
      headers: { origin: 'http://localhost:5173' },
    }),
  );

  assertEquals(res.status, 204);
  assertEquals(
    res.headers.get('access-control-allow-origin'),
    'http://localhost:5173',
  );
  assertEquals(
    res.headers.get('access-control-allow-methods'),
    'GET,POST,PATCH,DELETE,OPTIONS',
  );
});
