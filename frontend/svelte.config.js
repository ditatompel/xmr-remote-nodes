import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import * as child_process from 'node:child_process';
import { readFileSync } from 'fs';
import { join } from 'path';

// Helper function to execute shell commands
function execSync(cmd) {
	return child_process.execSync(cmd).toString().trim();
}

// Read version from package.json
const packageJsonPath = join(process.cwd(), 'package.json');
const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf8'));
const VERSION = packageJson.version;

// Retrieve current branch
const BRANCH = execSync('git rev-parse --abbrev-ref HEAD');

// Retrieve current tag if it exists
const RELEASE_TAG = execSync('git tag -l --points-at HEAD');

// Generate version suffix
const commitCount = execSync('git rev-list --count HEAD');
const shortCommitHash = execSync('git show --no-patch --no-notes --pretty="%h" HEAD');
const VERSION_SUFFIX = `-beta.${commitCount}.${shortCommitHash}`;

// Determine branch-specific values
let TAG_BRANCH = `.${BRANCH}`;

if (BRANCH === 'HEAD' || BRANCH === 'main') {
	TAG_BRANCH = '';
}

// Determine final tag
let TAG = `${VERSION}${VERSION_SUFFIX}${TAG_BRANCH}`;
if (RELEASE_TAG) {
	TAG = RELEASE_TAG;
}

console.log('Building with tag', TAG);

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		version: {
			name: TAG
		},
		// paths: {
		// 	base: '/'
		// },
		// trailingSlash: 'always',
		adapter: adapter()
	}
};

export default config;
