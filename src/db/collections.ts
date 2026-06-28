import { ObjectId } from 'mongodb';
import { db } from './mongo.ts';

export type UserDoc = {
  _id: ObjectId;
  username: string;
  passwordHash: string;
  createdAt: Date;
  updatedAt: Date;
};

export type RefreshTokenDoc = {
  _id: ObjectId;
  userId: ObjectId;
  tokenLookupHash: string;
  tokenHash: string;
  expiresAt: Date;
  revokedAt?: Date;
  createdAt: Date;
};

export type ProjectDoc = {
  _id: ObjectId;
  title: string;
  description: string;
  legend: {
    difficulty: Record<string, string>;
    type: Record<string, string>;
  };
  phases: Array<{
    id: string;
    title: string;
    type: string;
    difficulty: number;
    summary: string;
    conceptsToKnow: string[];
    toolsToLearn: string[];
    practiceLabs: string[];
    masteryChecks: string[];
    prerequisites: string[];
  }>;
  capstones: Array<{
    id: string;
    title: string;
    difficulty: number;
    summary: string;
    build: string[];
    conceptsToKnow: string;
    toolsToLearn: string[];
    prerequisites: string[];
  }>;
  recommendedOrder: string[];
  masteryDefinitions: string[];
  createdAt: Date;
  updatedAt: Date;
};

export type ProjectProgressDoc = {
  _id: ObjectId;
  userId: ObjectId;
  projectId: ObjectId;
  notes: Record<
    string,
    {
      notes: string[];
      links: { text: string; url: string }[];
    }
  >;
  createdAt: Date;
  updatedAt: Date;
};

// Temporary
export type GroupDoc = {
  _id: ObjectId;
  creatorId: ObjectId;
  name: string;
  people: string[];
  createdAt: Date;
  updatedAt: Date;
};

export type AiMessageDoc = {
  _id: ObjectId;
  userId: ObjectId;
  role: 'user' | 'assistant' | 'system';
  content: string;
  createdAt: Date;
};

export type AiExtractionDoc = {
  _id: ObjectId;
  userId: ObjectId;
  fileName?: string;
  mimeType: string;
  extractedText?: string;
  extractedObject?: unknown;
  createdAt: Date;
};

export const users = db.collection<UserDoc>('users');
export const refreshTokens = db.collection<RefreshTokenDoc>('refresh_tokens');
export const projects = db.collection<ProjectDoc>('projects');
export const projectProgresses =
  db.collection<ProjectProgressDoc>('projectProgresses');
export const groups = db.collection<GroupDoc>('groups');
export const aiMessages = db.collection<AiMessageDoc>('ai_messages');
export const aiExtractions = db.collection<AiExtractionDoc>('ai_extractions');

export async function ensureIndexes() {
  await users.createIndex({ username: 1 }, { unique: true });

  await refreshTokens.createIndex({ tokenLookupHash: 1 }, { unique: true });
  await refreshTokens.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 });
  await refreshTokens.createIndex({ userId: 1 });

  await aiExtractions.createIndex({ userId: 1, createdAt: -1 });
}
