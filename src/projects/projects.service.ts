import { InsertOneResult, ObjectId, WithId } from 'mongodb';
import {
  ProjectDoc,
  ProjectProgressDoc,
  projectProgresses,
  projects,
} from '../db/collections.ts';
import { HttpError } from '../shared/http.ts';
import { toObjectId } from '../shared/object-id.ts';

export type ProjectPayload = Omit<
  ProjectDoc,
  '_id' | 'createdAt' | 'updatedAt'
>;

export type ProjectProgressPayload = Omit<
  ProjectProgressDoc,
  '_id' | 'userId' | 'createdAt' | 'updatedAt'
>;

export async function listProjects(): Promise<WithId<ProjectDoc>[]> {
  return await projects.find().toArray();
}

export async function createProject(
  payload: ProjectPayload,
): Promise<InsertOneResult<ProjectDoc>> {
  const now = new Date();

  return await projects.insertOne({
    ...payload,
    _id: new ObjectId(),
    createdAt: now,
    updatedAt: now,
  });
}

export async function getProject(
  projectId: string,
): Promise<WithId<ProjectDoc>> {
  const project = await projects.findOne({
    _id: toObjectId(projectId),
  });

  if (!project) {
    throw new HttpError('Project not found', 404);
  }

  return project;
}

export async function updateProject(
  projectId: string,
  payload: ProjectPayload,
): Promise<WithId<ProjectDoc>> {
  const result = await projects.findOneAndUpdate(
    {
      _id: toObjectId(projectId),
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

export async function deleteProject(projectId: string): Promise<boolean> {
  const result = await projects.deleteOne({
    _id: toObjectId(projectId),
  });

  if (!result.deletedCount) {
    throw new HttpError('Project not found', 404);
  }

  return result.deletedCount === 1;
}

export async function startProject(
  userId: string,
  payload: ProjectProgressPayload,
): Promise<InsertOneResult<ProjectProgressDoc>> {
  const now = new Date();

  return await projectProgresses.insertOne({
    ...payload,
    _id: new ObjectId(),
    userId: toObjectId(userId),
    createdAt: now,
    updatedAt: now,
  });
}

export async function updateProjectProgress(
  userId: string,
  projectProgressId: string,
  payload: ProjectProgressPayload,
): Promise<WithId<ProjectProgressDoc>> {
  const result = await projectProgresses.findOneAndUpdate(
    {
      _id: toObjectId(projectProgressId),
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
    throw new HttpError('Project progress not found', 404);
  }

  return result;
}
