import { assertEquals, assertThrows } from '@std/assert';
import { validateAuthInput } from '../../src/auth/auth-validation.ts';
import { HttpError } from '../../src/shared/http.ts';

Deno.test('validateAuthInput normalizes email', () => {
  const result = validateAuthInput({
    email: '  USER@Example.COM  ',
    password: 'password123',
  });

  assertEquals(result, {
    email: 'user@example.com',
    password: 'password123',
  });
});

Deno.test('validateAuthInput rejects invalid email', () => {
  const error = assertThrows(
    () => validateAuthInput({ email: 'invalid', password: 'password123' }),
    HttpError,
  );

  assertEquals(error.status, 400);
});

Deno.test('validateAuthInput rejects short password', () => {
  const error = assertThrows(
    () => validateAuthInput({ email: 'user@example.com', password: 'short' }),
    HttpError,
  );

  assertEquals(error.status, 400);
});
