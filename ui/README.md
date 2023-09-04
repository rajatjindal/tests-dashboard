# Nuxt 3 Minimal Starter

Look at the [Nuxt 3 documentation](https://nuxt.com/docs/getting-started/introduction) to learn more.

## Setup

Make sure to install the dependencies:

```bash
# npm
npm install

# pnpm
pnpm install

# yarn
yarn install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
# npm
npm run dev

# pnpm
pnpm run dev

# yarn
yarn dev
```

## Production

Build the application for production:

```bash
# npm
npm run build

# pnpm
pnpm run build

# yarn
yarn build
```

Locally preview production build:

```bash
# npm
npm run preview

# pnpm
pnpm run preview

# yarn
yarn preview
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.


fix typescript error:
```
// this one works; edit: ./node_modules/nitropack/dist/h3-32ec41d1.d.ts
type MatchedRoutes<Route extends string, MatchedKeysResult extends MatchResult<string> = MatchResult<keyof InternalApi>, Matches extends MatchResult<string> = _MatchedRoutes<Route, MatchedKeysResult>> = Route extends "/" ? keyof InternalApi : Extract<Matches, {
    exact: true;
}> extends never ? Extract<Matches, { catchAll: true; }>["key"] : Extract<Matches, { exact: true; }>["key"];
```