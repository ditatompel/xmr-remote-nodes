# UI

Everything you need to build a Svelte project, powered by [`create-svelte`](https://github.com/sveltejs/kit/tree/main/packages/create-svelte).

## Developing

Once you've created a project and installed dependencies with `npm install` (or `pnpm install` or `yarn`), start a development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://kit.svelte.dev/docs/adapters) for your target environment.

## Deploying

after running `npm run build` from development device, copy `./build`, `package.json` and `package-lock.json` to server. On the server, run `npm ci --omit dev` then restart the systemd service.

Playbook example (run from root project):
```shell
ansible-playbook -i ./utils/ansible/inventory.ini -l production ./utils/ansible/deploy.yml -K
```
