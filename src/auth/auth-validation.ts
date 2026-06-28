import { HttpError } from '../shared/http.ts';

type AuthInput = {
  username?: unknown;
  password?: unknown;
};

export function validateAuthInput(data: AuthInput) {
  if (typeof data.username !== 'string') {
    throw new HttpError('Valid username is required', 400);
  }

  if (typeof data.password !== 'string' || data.password.length < 8) {
    throw new HttpError('Password must be at least 8 characters', 400);
  }

  return {
    username: data.username.trim().toLowerCase(),
    password: data.password,
  };
}
