import { requireUser } from '../auth/requireUser.ts';
import { Route } from '../routes.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import {
  listProjects,
  createProject,
  ProjectPayload,
  getProject,
} from './projects.service.ts';

export const projectRoutes: Route[] = [
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/projects' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const projects = await listProjects(user.userId);

      return json({ projects });
    },
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/projects/create' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const payload = await readJson<ProjectPayload>(req);
      const project = await createProject(user.userId, payload);

      return json({ project }, { status: 201 });
    },
  },
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/projects/:id' }),
    handler: async (req: Request, params: Record<string, string>) => {
      const user = await requireUser(req);
      const project = await getProject(user.userId, params.id);

      return json({ project });
    },
  },
];
