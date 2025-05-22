import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const response = await fetch('http://gradewise-api-backend:8080/');
	const data = await response.text();
	return { message: data };
};
