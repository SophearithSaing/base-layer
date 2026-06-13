import { InsertOneResult, ObjectId, WithId } from 'mongodb';
import { ProjectDoc, projects } from '../db/collections.ts';
import { HttpError } from '../shared/http.ts';
import { toObjectId } from '../shared/objectId.ts';

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
    userId: toObjectId(userId),
    name: payload.name,
    description: payload.description,
    createdAt: now,
    updatedAt: now,
  });
}

export async function getProject(
  userId: string,
  projectId: string,
): Promise<WithId<ProjectDoc>> {
  const project = await projects.findOne({
    _id: toObjectId(projectId),
    userId: toObjectId(userId),
  });

  if (!project) {
    throw new HttpError('Project not found', 404);
  }

  return project;
}

export async function updateProject(
  userId: string,
  projectId: string,
  payload: ProjectPayload,
): Promise<WithId<ProjectDoc>> {
  const result = await projects.findOneAndUpdate(
    {
      _id: toObjectId(projectId),
      userId: toObjectId(userId),
    },
    {
      $set: { ...payload, updatedAt: new Date() },
    },
    {
      returnDocument: 'after',
    },
  );

  if (!result) {
    throw new HttpError('Project not found', 404);
  }

  return result;
}

export async function deleteProject(
  userId: string,
  projectId: string,
): Promise<boolean> {
  const result = await projects.deleteOne({
    _id: toObjectId(projectId),
    userId: toObjectId(userId),
  });

  if (!result.deletedCount) {
    throw new HttpError('Project not found', 404);
  }

  return result.deletedCount === 1;
}
