import { HttpError } from '../shared/http.ts';

type AuthInput = {
  email?: unknown;
  password?: unknown;
};

export function validateAuthInput(data: AuthInput) {
  if (typeof data.email !== 'string' || !data.email.includes('@')) {
    throw new HttpError('Valid email is required', 400);
  }

  if (typeof data.password !== 'string' || data.password.length < 8) {
    throw new HttpError('Password must be at least 8 characters', 400);
  }

  return {
    email: data.email.trim().toLowerCase(),
    password: data.password,
  };
}
