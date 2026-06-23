import { assertEquals, assertStringIncludes } from '@std/assert';
import {
  clearCookie,
  getCookie,
  serializeCookie,
} from '../../src/shared/cookies.ts';

Deno.test('serializeCookie writes attribute name', () => {
  const cookie = serializeCookie('access_token', 'abc', {
    httpOnly: true,
    path: '/',
    maxAge: 60,
  });

  assertStringIncludes(cookie, 'HttpOnly');
  assertStringIncludes(cookie, 'SameSite=Lax');
  assertStringIncludes(cookie, 'Path=/');
  assertStringIncludes(cookie, 'Max-Age=60');
});

Deno.test('clearCookie expires cookie at the requested path', () => {
  const cookie = clearCookie('refresh_token', '/auth/refresh');

  assertStringIncludes(cookie, 'refresh_token=');
  assertStringIncludes(cookie, 'Path=/auth/refresh');
  assertStringIncludes(cookie, 'Max-Age=0');
});

Deno.test('getCookie reads a cookie value by name', () => {
  const req = new Request('http://localhost:8000/me', {
    headers: {
      cookie: 'access_token=abc; refresh_token=def',
    },
  });

  assertEquals(getCookie(req, 'refresh_token'), 'def');
});
