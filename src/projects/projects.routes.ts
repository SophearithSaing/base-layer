import { requireUser } from '../auth/require-user.ts';
import { Route } from '../routes.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import {
  createProject,
  deleteProject,
  getProject,
  listProjects,
  ProjectPayload,
  updateProject,
} from './projects.service.ts';
import {
  validateCreateProject,
  validateUpdateProject,
} from './project-validation.ts';

export const projectRoutes: Route[] = [
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/projects' }),
    handler: async (req: Request) => {
      await requireUser(req);
      const projects = await listProjects();

      return json({ projects });
    },
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/projects/create' }),
    handler: async (req: Request) => {
      await requireUser(req);
      const payload = await readJson<ProjectPayload>(req);
      validateCreateProject(payload);
      const project = await createProject(payload);

      return json({ project }, { status: 201 });
    },
  },
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/projects/:id' }),
    handler: async (req: Request, params: Record<string, string>) => {
      await requireUser(req);
      const project = await getProject(params.id);

      return json({ project });
    },
  },
  {
    method: HttpMethod.PATCH,
    pattern: new URLPattern({ pathname: '/projects/:id' }),
    handler: async (req: Request, params: Record<string, string>) => {
      await requireUser(req);
      const payload = await readJson<ProjectPayload>(req);
      validateUpdateProject(payload);
      const project = await updateProject(params.id, payload);

      return json({ project });
    },
  },
  {
    method: HttpMethod.DELETE,
    pattern: new URLPattern({ pathname: '/projects/:id' }),
    handler: async (req: Request, params: Record<string, string>) => {
      await requireUser(req);
      await deleteProject(params.id);

      return json({ ok: true });
    },
  },
];
