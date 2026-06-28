import { HttpError } from '../shared/http.ts';
import { ProjectPayload } from './projects.service.ts';

export function validateCreateProject(data: ProjectPayload) {
  if (typeof data.title !== 'string' || data.title.trim().length === 0) {
    throw new HttpError('Project title is required', 400);
  }
}

export function validateUpdateProject(data: ProjectPayload) {
  if (data.title !== undefined) {
    if (typeof data.title !== 'string' || data.title.trim().length === 0) {
      throw new HttpError('Project title cannot be empty', 400);
    }
  }

  if (data.description !== undefined) {
    if (typeof data.description !== 'string') {
      throw new HttpError('Project description must be a string', 400);
    }
  }
}
