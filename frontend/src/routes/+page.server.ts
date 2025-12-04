import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';

export const load: PageServerLoad = async () => {
  const base = env.API_URL ?? 'http://gradewise-api-backend';
  try {
    // TODO: replace '/health' with your real endpoint
    const r = await fetch(`${base}/api/health`);
    const text = await r.text();
    return { message: text };
  } catch {
    return { message: 'api-unavailable' }; // do not throw; avoid 500
  }
};