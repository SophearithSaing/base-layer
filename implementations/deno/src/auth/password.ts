import { hash, verify } from '@felix/argon2';

export async function hashPassword(password: string): Promise<string> {
  return await hash(password);
}

export async function verifyPassword(
  password: string,
  passwordHash: string,
): Promise<boolean> {
  return await verify(passwordHash, password);
}
