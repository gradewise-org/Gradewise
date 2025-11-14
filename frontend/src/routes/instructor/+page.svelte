<script lang="ts">
  import { PUBLIC_BASE_URL } from "$env/static/public";
  import FileDropzone from "$lib/FileDropzone.svelte";
  import { base } from '$app/paths';


  let assignmentName = "";
  let graderFile: File | null = null;
  let loading = false;
  let createdAssignmentId: string | null = null;
  let studentSubmitUrl: string | null = null;
  let error: string | null = null;

  function handleFilesSelected(event: CustomEvent<{ files: FileList }>) {
    graderFile = event.detail.files[0]; // just first file
  }

  async function createAssignment() {
    error = null;

    if (!assignmentName || !graderFile) {
      error = "Please provide an assignment name and grader file.";
      return;
    }

    loading = true;

    try {
      const formData = new FormData();
      formData.append("name", assignmentName);
      formData.append("grader", graderFile);

      const res = await fetch(`${PUBLIC_BASE_URL}/api/assignments`, {
        method: "POST",
        body: formData
      });

      if (!res.ok) {
        throw new Error(`Failed to create assignment (${res.status})`);
      }

      const data = await res.json();
      createdAssignmentId = data.assignmentId;
      studentSubmitUrl = data.studentSubmitUrl ?? `${base}/submit/${data.assignmentId}`;
    } catch (e: any) {
      error = e?.message ?? "Something went wrong.";
    } finally {
      loading = false;
    }
  }
</script>

<h1 class="text-2xl font-bold mb-4">Instructor: Create Assignment</h1>

<div class="space-y-4">
  <div class="space-y-2">
    <label class="block text-sm font-medium">Assignment Name</label>
    <input
      class="w-full border rounded px-3 py-2"
      bind:value={assignmentName}
      placeholder="e.g. Project 1 - Cache Simulator"
    />
  </div>

  <div class="space-y-2">
    <label class="block text-sm font-medium">Upload Grader (zip/tgz)</label>
    <FileDropzone on:filesSelected={handleFilesSelected} />
    {#if graderFile}
      <p class="text-xs text-gray-500 mt-1">
        Selected: {graderFile.name}
      </p>
    {/if}
  </div>

  {#if error}
    <p class="text-sm text-red-600">{error}</p>
  {/if}

  <button
    class="px-4 py-2 rounded bg-blue-600 text-white disabled:opacity-60"
    on:click={createAssignment}
    disabled={loading}
  >
    {#if loading}
      Creating...
    {:else}
      Create Assignment
    {/if}
  </button>

  {#if createdAssignmentId}
    <div class="mt-6 p-4 border rounded bg-green-50 text-sm space-y-1">
      <p><strong>Assignment ID:</strong> {createdAssignmentId}</p>
      {#if studentSubmitUrl}
        <p>
          <strong>Student Submit URL:</strong>
          <a href={studentSubmitUrl} class="text-blue-600 underline">
            {studentSubmitUrl}
          </a>
        </p>
      {/if}
    </div>
  {/if}
</div>