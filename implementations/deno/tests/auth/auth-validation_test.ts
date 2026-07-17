import { assertEquals, assertThrows } from '@std/assert';
import { validateAuthInput } from '../../src/auth/auth-validation.ts';
import { HttpError } from '../../src/shared/http.ts';

Deno.test('validateAuthInput normalizes username', () => {
  const result = validateAuthInput({
    username: '  USER_Example  ',
    password: 'password123',
  });

  assertEquals(result, {
    username: 'user_example',
    password: 'password123',
  });
});

Deno.test('validateAuthInput rejects short password', () => {
  const error = assertThrows(
    () => validateAuthInput({ username: 'user', password: 'short' }),
    HttpError,
  );

  assertEquals(error.status, 400);
});
