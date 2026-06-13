import { InsertOneResult, ObjectId, WithId } from 'mongodb';
import { ProjectDoc, projects } from '../db/collections.ts';

export type ProjectPayload = {
  name: string;
  description: string;
};

export async function listProjects(
  userId: string,
): Promise<WithId<ProjectDoc>[]> {
  return await projects.find({ userId: new ObjectId(userId) }).toArray();
}

export async function createProject(
  userId: string,
  payload: ProjectPayload,
): Promise<InsertOneResult<ProjectDoc>> {
  const now = new Date();

  return await projects.insertOne({
    _id: new ObjectId(),
    userId: new ObjectId(userId),
    name: payload.name,
    description: payload.description,
    createdAt: now,
    updatedAt: now,
  });
}
