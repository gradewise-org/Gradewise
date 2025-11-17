<script lang="ts">
    import { onMount } from 'svelte';
    import { get } from 'svelte/store';
    import { ltiSession } from '$lib/stores/ltiSession';

    let code = '';
    let term = '';
    let title = '';
    let section = '';
    let startDate = '';
    let endDate = '';

    let submitting = false;
    let successMessage: string | null = null;
    let errorMessage: string | null = null;

    // Optional: default term (you can change/remove this)
    onMount(() => {
        if (!term) term = 'Fall 2025';
    });

    async function handleSubmit(event: SubmitEvent) {
        event.preventDefault();
        successMessage = null;
        errorMessage = null;
        submitting = true;

        try {
            const { facultyId } = get(ltiSession);

            if (!facultyId) {
                errorMessage = 'No facultyId found from LTI session. Try relaunching from Canvas.';
                return;
            }

            if (!code || !term || !title || !section) {
                errorMessage = 'Please fill in course code, term, title, and section.';
                return;
            }

            const hasStart = !!startDate;
            const hasEnd = !!endDate;

            if (hasStart !== hasEnd) {
                errorMessage = 'If you provide a start date, you must also provide an end date (and vice versa).';
                return;
            }

            const payload: any = {
                code,
                term,
                title,
                section,
            };

            if (hasStart && hasEnd) {
                payload.startDate = startDate;
                payload.endDate = endDate;
            }

            const res = await fetch('/api/courses', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Faculty-ID': facultyId
                },
                body: JSON.stringify(payload),
                credentials: 'include'
            });

            const text = await res.text();
            let data: any = null;
            try {
                data = text ? JSON.parse(text) : null;
            } catch {
                // not JSON, leave as text
            }

            if (res.status === 201) {
                successMessage = 'Course created successfully.';
                errorMessage = null;

            } else {
                const apiMessage =
                    data && typeof data === 'object'
                        ? (data.message || data.error || JSON.stringify(data))
                        : text || 'Unknown error';

                errorMessage = `Failed to create course: ${res.status} ${res.statusText} – ${apiMessage}`;
            }
        } catch (err: any) {
            console.error(err);
            errorMessage = `Unexpected error: ${err?.message ?? String(err)}`;
        } finally {
            submitting = false;
        }
    }
</script>

<svelte:head>
    <title>Create Course • Gradewise</title>
</svelte:head>

<div class="space-y-6">
    <h1 class="text-2xl font-semibold text-slate-900">Create a new course</h1>
    <p class="text-sm text-slate-600">
        Fill in the course details below. This will call <code>POST /api/courses</code> on the backend.
    </p>

    {#if successMessage}
        <div class="rounded-md border border-green-200 bg-green-50 px-3 py-2 text-sm text-green-800">
            {successMessage}
        </div>
    {/if}

    {#if errorMessage}
        <div class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">
            {errorMessage}
        </div>
    {/if}

    <form class="space-y-4 max-w-xl" on:submit|preventDefault={handleSubmit}>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
                <label class="block text-sm font-medium text-slate-700">
                    Course code
                </label>
                <input
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={code}
                        placeholder="CS101"
                        required
                />
            </div>

            <div>
                <label class="block text-sm font-medium text-slate-700">
                    Section
                </label>
                <input
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={section}
                        placeholder="001"
                        required
                />
            </div>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
                <label class="block text-sm font-medium text-slate-700">
                    Term
                </label>
                <input
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={term}
                        placeholder="Fall 2025"
                        required
                />
            </div>

            <div>
                <label class="block text-sm font-medium text-slate-700">
                    Title
                </label>
                <input
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={title}
                        placeholder="Intro to CS"
                        required
                />
            </div>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
                <label class="block text-sm font-medium text-slate-700">
                    Start date
                </label>
                <input
                        type="date"
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={startDate}
                />
            </div>

            <div>
                <label class="block text-sm font-medium text-slate-700">
                    End date
                </label>
                <input
                        type="date"
                        class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                        bind:value={endDate}
                />
            </div>
        </div>

        <div class="pt-2">
            <button
                    type="submit"
                    class="inline-flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
                    disabled={submitting}
            >
                {#if submitting}
                    Creating…
                {:else}
                    Create course
                {/if}
            </button>
        </div>
    </form>
</div>
