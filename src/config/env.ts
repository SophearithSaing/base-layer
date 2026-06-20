import '@std/dotenv/load';

function getRequiredEnv(key: string): string {
  const value = Deno.env.get(key);

  if (!value) {
    throw new Error(`Missing required env var: ${key}`);
  }

  return value;
}

export const env = {
  APP_ENV: Deno.env.get('APP_ENV') ?? 'dev',
  PORT: Number(Deno.env.get('PORT')) ?? 8000,
  CLIENT_ORIGIN: Deno.env.get('CLIENT_ORIGIN') ?? 'http://localhost:5173',

  MONGODB_URI: getRequiredEnv('MONGODB_URI'),
  MONGODB_DB_NAME: Deno.env.get('MONGODB_DB_NAME') ?? 'baselayer',

  JWT_SECRET: getRequiredEnv('JWT_SECRET'),
  TOGETHER_API_KEY: Deno.env.get('TOGETHER_API_KEY'),
};
