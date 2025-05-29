import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const response = await fetch('http://gradewise-api-backend/');
	const data = await response.text();
	return { message: data };
};
