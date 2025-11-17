<script lang="ts">
    import type { PageData } from './$types';

    export let data: PageData;

    // Bind local inputs to keep the form controlled
    let term = data.term ?? '';
    let code = data.code ?? '';
</script>

<svelte:head>
    <title>My Courses • Gradewise</title>
</svelte:head>

<div class="space-y-6">
    <h1 class="text-2xl font-semibold text-slate-900">My courses</h1>
    <p class="text-sm text-slate-600">
        These are your courses (using <code>GET /api/courses?mine=true</code>).
        You can optionally filter by term and course code.
    </p>

    {#if data.error}
        <div class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-800">
            {data.error}
        </div>
    {/if}

    <!-- Filter form: GET so filters appear in URL (?term=...&code=...) -->
    <form
            class="flex flex-wrap items-end gap-4 border-b border-slate-200 pb-4 mb-4"
            method="GET"
    >
        <div>
            <label class="block text-sm font-medium text-slate-700">
                Term (optional)
            </label>
            <input
                    class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    name="term"
                    bind:value={term}
                    placeholder="Fall 2025"
            />
        </div>

        <div>
            <label class="block text-sm font-medium text-slate-700">
                Course code (optional)
            </label>
            <input
                    class="mt-1 w-full rounded-md border border-slate-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                    name="code"
                    bind:value={code}
                    placeholder="CS101"
            />
        </div>

        <button
                type="submit"
                class="inline-flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700"
        >
            Apply filters
        </button>
    </form>

    <!-- Courses table -->
    {#if data.courses.length === 0}
        <p class="text-sm text-slate-600">
            No courses found for these filters.
        </p>
    {:else}
        <div class="overflow-x-auto rounded-lg border border-slate-200">
            <table class="min-w-full text-sm">
                <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500">
                <tr>
                    <th class="px-3 py-2">Code</th>
                    <th class="px-3 py-2">Term</th>
                    <th class="px-3 py-2">Title</th>
                    <th class="px-3 py-2">Section</th>
                    <th class="px-3 py-2">Start</th>
                    <th class="px-3 py-2">End</th>
                    <th class="px-3 py-2">Created</th>
                </tr>
                </thead>
                <tbody class="divide-y divide-slate-200 bg-white">
                {#each data.courses as c}
                    <tr>
                        <td class="px-3 py-2 font-mono">{c.code}</td>
                        <td class="px-3 py-2">{c.term}</td>
                        <td class="px-3 py-2">{c.title}</td>
                        <td class="px-3 py-2">{c.section ?? '—'}</td>
                        <td class="px-3 py-2">{c.startDate ?? '—'}</td>
                        <td class="px-3 py-2">{c.endDate ?? '—'}</td>
                        <td class="px-3 py-2 text-xs text-slate-500">
                            {new Date(c.createdAt).toLocaleString()}
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>
