import { writable } from 'svelte/store';

type LTIUser = {
    name?: string | null;
    email?: string | null;
};

type LTISession = {
    facultyId: string | null;
    user?: LTIUser;
};
export const ltiSession = writable<LTISession>({ facultyId: null });