// src/routes/instructor/courses/+page.server.ts
import type { PageServerLoad } from './$types';
import { env } from '$env/dynamic/private';

const API_BASE = env.API_URL ?? 'http://gradewise-api-backend';

type Course = {
    id: string;
    code: string;
    term: string;
    title: string;
    section: string | null;
    startDate: string | null;
    endDate: string | null;
    creatorFacultyId: string;
    createdAt: string;
};

export const load: PageServerLoad = async ({ parent, url, fetch }) => {
    // Get LTI data from the layout.server.ts load
    const { lti } = await parent();
    const facultyId = lti?.facultyId;

    if (!facultyId) {
        return {
            courses: [] as Course[],
            term: '',
            code: '',
            error: 'No facultyId from LTI. Try relaunching from Canvas.',
        };
    }

    const term = url.searchParams.get('term') ?? '';
    const code = url.searchParams.get('code') ?? '';

    const params = new URLSearchParams();
    params.set('mine', 'true');
    if (term) params.set('term', term);
    if (code) params.set('code', code);

    let courses: Course[] = [];
    let error: string | null = null;

    try {
        const res = await fetch(`${API_BASE}/api/courses?${params.toString()}`, {
            headers: {
                'X-Faculty-ID': facultyId,
            },
        });

        if (res.ok) {
            courses = await res.json();
        } else {
            const text = await res.text();
            error = `Failed to load courses: ${res.status} ${text || res.statusText}`;
        }
    } catch (e: any) {
        error = `Unexpected error loading courses: ${e?.message ?? String(e)}`;
    }

    return {
        courses,
        term,
        code,
        error,
    };
};
