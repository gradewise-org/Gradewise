import type { LayoutServerLoad } from './$types';
import { env } from '$env/dynamic/private';

type LTIUser = {
    name?: string | null;
    email?: string | null;
};

type LTISession = {
    facultyId: string | null;
    user?: LTIUser;
};
const lti_base = env.LTI_BASE ?? 'http://gradewise-lti';
export const load: LayoutServerLoad = async ({ fetch, cookies }) => {
    const launched = cookies.get('launched') === '1';

    let lti: LTISession = { facultyId: null };

    try {
        const sid = cookies.get('__Host-gw_sid');
        const res = await fetch(`${lti_base}/lti/session`, {
            // IMPORTANT: SvelteKit's server fetch does NOT automatically forward
            // cookies to external services. You must send it explicitly:
            headers: sid
                ? {
                    cookie: `__Host-gw_sid=${sid}`,
                }
                : {},
        });
        console.log('layout.server.ts /lti/session status:', res.status);
        if (res.ok) {
            const data = await res.json();
            lti = {
                facultyId: data.facultyId ?? null,
                user: data.user ?? {},
            };
        } else {
            // Optional: log or ignore
            console.error('lti/session failed in layout.server.ts', res.status, await res.text());
        }
    } catch (e) {
        console.error('Error calling /lti/session from layout.server.ts', e);
    }

    return {
        launched,
        lti,
    };
};