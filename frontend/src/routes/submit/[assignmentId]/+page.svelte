<script lang="ts">
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import { page } from "$app/stores";
  import FileDropzone from "$lib/FileDropzone.svelte";

  let assignmentId = "";
  let submissionFile: File | null = null;
  let loading = false;
  let error: string | null = null;
  let submissionId: string | null = null;
  let status: string | null = null;

  $: assignmentId = $page.params.assignmentId;

  function handleFilesSelected(event: CustomEvent<{ files: FileList }>) {
    submissionFile = event.detail.files[0];
  }

  async function submit() {
    error = null;
    if (!submissionFile) {
      error = "Please select a submission file.";
      return;
    }

    loading = true;

    try {
      const formData = new FormData();
      formData.append("submission", submissionFile);

      const res = await fetch(
        `${PUBLIC_BASE_URL}/api/assignments/${assignmentId}/submissions`,
        {
          method: "POST",
          body: formData
        }
      );

      if (!res.ok) throw new Error(`Failed to submit (${res.status})`);

      const data = await res.json();
      submissionId = data.submissionId;
      status = data.status ?? "queued";
    } catch (e: any) {
      error = e?.message ?? "Something went wrong.";
    } finally {
      loading = false;
    }
  }

  async function refreshStatus() {
    if (!submissionId) return;

    const res = await fetch(
      `${PUBLIC_BASE_URL}/api/submissions/${submissionId}`
    );
    if (res.ok) {
      const data = await res.json();
      status = data.status;
      // Optionally: data.score, data.logsUrl, etc.
    }
  }
</script>

<h1 class="text-2xl font-bold mb-2">Submit Assignment</h1>
<p class="text-sm text-gray-600 mb-4">Assignment ID: {assignmentId}</p>

<div class="space-y-4">
  <div class="space-y-2">
    <label class="block text-sm font-medium">
      Upload your submission (zip/tgz)
    </label>
    <FileDropzone on:filesSelected={handleFilesSelected} />
    {#if submissionFile}
      <p class="text-xs text-gray-500 mt-1">
        Selected: {submissionFile.name}
      </p>
    {/if}
  </div>

  {#if error}
    <p class="text-sm text-red-600">{error}</p>
  {/if}

  <div class="flex items-center gap-3">
    <button
      class="px-4 py-2 rounded bg-blue-600 text-white disabled:opacity-60"
      on:click={submit}
      disabled={loading}
    >
      {#if loading}
        Submitting...
      {:else}
        Submit
      {/if}
    </button>

    {#if submissionId}
      <button
        class="px-3 py-2 rounded border text-sm"
        on:click={refreshStatus}
      >
        Refresh status
      </button>
    {/if}
  </div>

  {#if submissionId}
    <div class="mt-4 p-3 border rounded text-sm">
      <p><strong>Submission ID:</strong> {submissionId}</p>
      <p><strong>Status:</strong> {status}</p>
    </div>
  {/if}
</div>