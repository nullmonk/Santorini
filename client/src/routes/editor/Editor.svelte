<script lang="ts">
    import CodeMirror from "svelte-codemirror-editor";
    import { python } from "@codemirror/lang-python";

    let value = "";
</script>

<CodeMirror bind:value lang={python()} tabSize={2} />

<style>
    .textEdit {
        height: 100%;
    }
    .wrapper {
        display:flex;
        flex-grow: 1;
        height: 100%;
    }
</style>