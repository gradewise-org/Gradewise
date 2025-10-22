import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
  const launched = cookies.get('launched') === '1';
  return { launched };
};