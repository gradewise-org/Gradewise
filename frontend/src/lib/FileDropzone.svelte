<script lang="ts">
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
  let dragging = false;

  function handleDrop(event: DragEvent) {
    event.preventDefault();
    dragging = false;
    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      dispatch("filesSelected", { files });
    }
  }

  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      dispatch("filesSelected", { files: input.files });
    }
  }
</script>

<div
  class={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer ${
    dragging ? "border-blue-500 bg-blue-50" : "border-gray-300 bg-white"
  }`}
  on:dragover|preventDefault={() => (dragging = true)}
  on:dragleave={() => (dragging = false)}
  on:drop={handleDrop}
>
  <p class="mb-2 font-semibold">Drag &amp; drop your file here</p>
  <p class="text-sm text-gray-500 mb-4">or click to browse</p>
  <input
    type="file"
    class="hidden"
    id="file-input"
    on:change={handleFileChange}
  />
  <label for="file-input" class="underline text-blue-600 text-sm cursor-pointer">
    Choose file
  </label>
</div>