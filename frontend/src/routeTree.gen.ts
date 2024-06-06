/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from '@tanstack/react-router'

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as SignupImport } from './routes/signup'
import { Route as SignoutImport } from './routes/signout'
import { Route as SigninImport } from './routes/signin'
import { Route as RecoverImport } from './routes/recover'
import { Route as AuthRouteImport } from './routes/_auth/route'
import { Route as AuthWorkspacesIndexImport } from './routes/_auth/workspaces/index'
import { Route as AuthWorkspacesAddImport } from './routes/_auth/workspaces/add'
import { Route as AuthWorkspacesWorkspaceIdImport } from './routes/_auth/workspaces/$workspaceId'
import { Route as AuthWorkspacesWorkspaceIdWorkspaceRouteImport } from './routes/_auth/workspaces/$workspaceId/_workspace/route'
import { Route as AuthWorkspacesWorkspaceIdSetupIndexImport } from './routes/_auth/workspaces/$workspaceId/setup/index'
import { Route as AuthWorkspacesWorkspaceIdWorkspaceIndexImport } from './routes/_auth/workspaces/$workspaceId/_workspace/index'
import { Route as AuthWorkspacesWorkspaceIdThreadsThreadIdImport } from './routes/_auth/workspaces/$workspaceId/threads/$threadId'
import { Route as AuthWorkspacesWorkspaceIdWorkspaceUnassignedImport } from './routes/_auth/workspaces/$workspaceId/_workspace/unassigned'
import { Route as AuthWorkspacesWorkspaceIdWorkspaceMeImport } from './routes/_auth/workspaces/$workspaceId/_workspace/me'
import { Route as AuthWorkspacesWorkspaceIdWorkspaceLabelsLabelIdImport } from './routes/_auth/workspaces/$workspaceId/_workspace/labels.$labelId'

// Create Virtual Routes

const AuthIndexLazyImport = createFileRoute('/_auth/')()

// Create/Update Routes

const SignupRoute = SignupImport.update({
  path: '/signup',
  getParentRoute: () => rootRoute,
} as any)

const SignoutRoute = SignoutImport.update({
  path: '/signout',
  getParentRoute: () => rootRoute,
} as any)

const SigninRoute = SigninImport.update({
  path: '/signin',
  getParentRoute: () => rootRoute,
} as any)

const RecoverRoute = RecoverImport.update({
  path: '/recover',
  getParentRoute: () => rootRoute,
} as any)

const AuthRouteRoute = AuthRouteImport.update({
  id: '/_auth',
  getParentRoute: () => rootRoute,
} as any)

const AuthIndexLazyRoute = AuthIndexLazyImport.update({
  path: '/',
  getParentRoute: () => AuthRouteRoute,
} as any).lazy(() => import('./routes/_auth/index.lazy').then((d) => d.Route))

const AuthWorkspacesIndexRoute = AuthWorkspacesIndexImport.update({
  path: '/workspaces/',
  getParentRoute: () => AuthRouteRoute,
} as any)

const AuthWorkspacesAddRoute = AuthWorkspacesAddImport.update({
  path: '/workspaces/add',
  getParentRoute: () => AuthRouteRoute,
} as any)

const AuthWorkspacesWorkspaceIdRoute = AuthWorkspacesWorkspaceIdImport.update({
  path: '/workspaces/$workspaceId',
  getParentRoute: () => AuthRouteRoute,
} as any)

const AuthWorkspacesWorkspaceIdWorkspaceRouteRoute =
  AuthWorkspacesWorkspaceIdWorkspaceRouteImport.update({
    id: '/_workspace',
    getParentRoute: () => AuthWorkspacesWorkspaceIdRoute,
  } as any)

const AuthWorkspacesWorkspaceIdSetupIndexRoute =
  AuthWorkspacesWorkspaceIdSetupIndexImport.update({
    path: '/setup/',
    getParentRoute: () => AuthWorkspacesWorkspaceIdRoute,
  } as any)

const AuthWorkspacesWorkspaceIdWorkspaceIndexRoute =
  AuthWorkspacesWorkspaceIdWorkspaceIndexImport.update({
    path: '/',
    getParentRoute: () => AuthWorkspacesWorkspaceIdWorkspaceRouteRoute,
  } as any)

const AuthWorkspacesWorkspaceIdThreadsThreadIdRoute =
  AuthWorkspacesWorkspaceIdThreadsThreadIdImport.update({
    path: '/threads/$threadId',
    getParentRoute: () => AuthWorkspacesWorkspaceIdRoute,
  } as any)

const AuthWorkspacesWorkspaceIdWorkspaceUnassignedRoute =
  AuthWorkspacesWorkspaceIdWorkspaceUnassignedImport.update({
    path: '/unassigned',
    getParentRoute: () => AuthWorkspacesWorkspaceIdWorkspaceRouteRoute,
  } as any)

const AuthWorkspacesWorkspaceIdWorkspaceMeRoute =
  AuthWorkspacesWorkspaceIdWorkspaceMeImport.update({
    path: '/me',
    getParentRoute: () => AuthWorkspacesWorkspaceIdWorkspaceRouteRoute,
  } as any)

const AuthWorkspacesWorkspaceIdWorkspaceLabelsLabelIdRoute =
  AuthWorkspacesWorkspaceIdWorkspaceLabelsLabelIdImport.update({
    path: '/labels/$labelId',
    getParentRoute: () => AuthWorkspacesWorkspaceIdWorkspaceRouteRoute,
  } as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/_auth': {
      id: '/_auth'
      path: ''
      fullPath: ''
      preLoaderRoute: typeof AuthRouteImport
      parentRoute: typeof rootRoute
    }
    '/recover': {
      id: '/recover'
      path: '/recover'
      fullPath: '/recover'
      preLoaderRoute: typeof RecoverImport
      parentRoute: typeof rootRoute
    }
    '/signin': {
      id: '/signin'
      path: '/signin'
      fullPath: '/signin'
      preLoaderRoute: typeof SigninImport
      parentRoute: typeof rootRoute
    }
    '/signout': {
      id: '/signout'
      path: '/signout'
      fullPath: '/signout'
      preLoaderRoute: typeof SignoutImport
      parentRoute: typeof rootRoute
    }
    '/signup': {
      id: '/signup'
      path: '/signup'
      fullPath: '/signup'
      preLoaderRoute: typeof SignupImport
      parentRoute: typeof rootRoute
    }
    '/_auth/': {
      id: '/_auth/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof AuthIndexLazyImport
      parentRoute: typeof AuthRouteImport
    }
    '/_auth/workspaces/$workspaceId': {
      id: '/_auth/workspaces/$workspaceId'
      path: '/workspaces/$workspaceId'
      fullPath: '/workspaces/$workspaceId'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdImport
      parentRoute: typeof AuthRouteImport
    }
    '/_auth/workspaces/add': {
      id: '/_auth/workspaces/add'
      path: '/workspaces/add'
      fullPath: '/workspaces/add'
      preLoaderRoute: typeof AuthWorkspacesAddImport
      parentRoute: typeof AuthRouteImport
    }
    '/_auth/workspaces/': {
      id: '/_auth/workspaces/'
      path: '/workspaces'
      fullPath: '/workspaces'
      preLoaderRoute: typeof AuthWorkspacesIndexImport
      parentRoute: typeof AuthRouteImport
    }
    '/_auth/workspaces/$workspaceId/_workspace': {
      id: '/_auth/workspaces/$workspaceId/_workspace'
      path: ''
      fullPath: '/workspaces/$workspaceId'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceRouteImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdImport
    }
    '/_auth/workspaces/$workspaceId/_workspace/me': {
      id: '/_auth/workspaces/$workspaceId/_workspace/me'
      path: '/me'
      fullPath: '/workspaces/$workspaceId/me'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceMeImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceRouteImport
    }
    '/_auth/workspaces/$workspaceId/_workspace/unassigned': {
      id: '/_auth/workspaces/$workspaceId/_workspace/unassigned'
      path: '/unassigned'
      fullPath: '/workspaces/$workspaceId/unassigned'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceUnassignedImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceRouteImport
    }
    '/_auth/workspaces/$workspaceId/threads/$threadId': {
      id: '/_auth/workspaces/$workspaceId/threads/$threadId'
      path: '/threads/$threadId'
      fullPath: '/workspaces/$workspaceId/threads/$threadId'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdThreadsThreadIdImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdImport
    }
    '/_auth/workspaces/$workspaceId/_workspace/': {
      id: '/_auth/workspaces/$workspaceId/_workspace/'
      path: '/'
      fullPath: '/workspaces/$workspaceId/'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceIndexImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceRouteImport
    }
    '/_auth/workspaces/$workspaceId/setup/': {
      id: '/_auth/workspaces/$workspaceId/setup/'
      path: '/setup'
      fullPath: '/workspaces/$workspaceId/setup'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdSetupIndexImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdImport
    }
    '/_auth/workspaces/$workspaceId/_workspace/labels/$labelId': {
      id: '/_auth/workspaces/$workspaceId/_workspace/labels/$labelId'
      path: '/labels/$labelId'
      fullPath: '/workspaces/$workspaceId/labels/$labelId'
      preLoaderRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceLabelsLabelIdImport
      parentRoute: typeof AuthWorkspacesWorkspaceIdWorkspaceRouteImport
    }
  }
}

// Create and export the route tree

export const routeTree = rootRoute.addChildren({
  AuthRouteRoute: AuthRouteRoute.addChildren({
    AuthIndexLazyRoute,
    AuthWorkspacesWorkspaceIdRoute: AuthWorkspacesWorkspaceIdRoute.addChildren({
      AuthWorkspacesWorkspaceIdWorkspaceRouteRoute:
        AuthWorkspacesWorkspaceIdWorkspaceRouteRoute.addChildren({
          AuthWorkspacesWorkspaceIdWorkspaceMeRoute,
          AuthWorkspacesWorkspaceIdWorkspaceUnassignedRoute,
          AuthWorkspacesWorkspaceIdWorkspaceIndexRoute,
          AuthWorkspacesWorkspaceIdWorkspaceLabelsLabelIdRoute,
        }),
      AuthWorkspacesWorkspaceIdThreadsThreadIdRoute,
      AuthWorkspacesWorkspaceIdSetupIndexRoute,
    }),
    AuthWorkspacesAddRoute,
    AuthWorkspacesIndexRoute,
  }),
  RecoverRoute,
  SigninRoute,
  SignoutRoute,
  SignupRoute,
})

/* prettier-ignore-end */
