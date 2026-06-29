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

export type ProjectProgressResponse = ProjectProgressDoc & {
  project: ProjectDoc;
};

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
  const project = await projects.findOne({
    _id: toObjectId(payload.projectId),
  });

  if (!project) {
    throw new HttpError('Project not found', 404);
  }

  return await projectProgresses.insertOne({
    ...payload,
    _id: new ObjectId(),
    userId: toObjectId(userId),
    title: project.title,
    description: project.description,
    progress: 0,
    createdAt: now,
    updatedAt: now,
  });
}

export async function listProjectProgresses(
  userId: string,
): Promise<ProjectProgressDoc[]> {
  return await projectProgresses.find({ userId: toObjectId(userId) }).toArray();
}

export async function getProjectProgress(
  userId: string,
  projectProgressId: string,
): Promise<ProjectProgressResponse> {
  const projectProgress = await projectProgresses.findOne({
    _id: toObjectId(projectProgressId),
    userId: toObjectId(userId),
  });

  if (!projectProgress) {
    throw new HttpError('Project progress not found', 404);
  }

  const project = await projects.findOne({
    _id: toObjectId(projectProgress.projectId),
  });

  if (!project) {
    throw new HttpError('Project not found', 404);
  }

  return { ...projectProgress, project };
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
