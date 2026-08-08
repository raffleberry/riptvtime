import { apiUploadImportZip } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { PAGE, theme } from "../../utils.js"
import { onMounted, ref } from "../../vue.js"

const Upload = {
  props: {},
  components: {},
  emits: ["success", "uploadDone", "uploadStart"],
  setup: (props, ctx) => {
    const isDragging = ref(false)
    const selectedFile = ref(null)
    const isUploading = ref(false)
    const uploadProgress = ref(0)

    const onDragOver = (e) => {
      if (isUploading.value) return
      e.preventDefault()
      isDragging.value = true
    }

    const onDragLeave = () => {
      isDragging.value = false
    }

    const onDrop = (e) => {
      if (isUploading.value) return
      e.preventDefault()
      isDragging.value = false

      if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
        selectedFile.value = e.dataTransfer.files[0]
      }
    }

    const onFileChange = (e) => {
      if (e.target.files && e.target.files.length > 0) {
        selectedFile.value = e.target.files[0]
      }
    }

    const triggerFileSelect = (fileInputRef) => {
      if (isUploading.value) return
      fileInputRef.click()
    }

    const startUpload = async () => {
      if (!selectedFile.value || isUploading.value) return

      isUploading.value = true
      uploadProgress.value = 20

      // Prepare multipart form data payload
      const formData = new FormData()
      formData.append("file", selectedFile.value)

      try {
        ctx.emit("uploadStart")
        const data = await apiUploadImportZip(formData, uploadProgress)
        uploadProgress.value = 100
        selectedFile.value = null
        ctx.emit("success", data)
      } catch (error) {
        console.error("Upload error details:", error)
        notify(MsgType.Error, "Upload", error)
      } finally {
        setTimeout(() => {
          isUploading.value = false
          ctx.emit("uploadDone")
          uploadProgress.value = 0
        }, 500)
      }
    }

    onMounted(() => {})

    return {
      isDragging,
      selectedFile,
      isUploading,
      uploadProgress,
      onDragOver,
      onDragLeave,
      onDrop,
      onFileChange,
      triggerFileSelect,
      startUpload,
    }
  },
  template: /* HTML */ `
    <h4 class="mt-4">Upload</h4>

    <div
      class="mb-3 card text-center p-5 border-2 position-relative"
      :class="[
                isDragging ? 'border-primary bg-primary bg-opacity-10' : 'border-secondary-subtle bg-body-tertiary',
                isUploading ? 'opacity-75' : ''
            ]"
      :style="{ 
                borderStyle: 'dashed !important', 
                cursor: isUploading ? 'not-allowed' : 'pointer' 
            }"
      @dragover="onDragOver"
      @dragleave="onDragLeave"
      @drop="onDrop"
      @click="triggerFileSelect($refs.fileInput)"
    >
      <div
        v-if="isUploading"
        class="position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-75 rounded"
        style="z-index: 10;"
      >
        <div
          class="spinner-border text-primary mb-3"
          role="status"
          style="width: 3rem; height: 3rem;"
        >
          <span class="visually-hidden">Loading...</span>
        </div>
        <h5 class="text-primary fw-bold mb-2">Uploading File...</h5>
        <div class="progress w-50" style="height: 6px;">
          <div
            class="progress-bar progress-bar-striped progress-bar-animated bg-primary"
            role="progressbar"
            :style="{ width: uploadProgress + '%' }"
          ></div>
        </div>
      </div>

      <div class="card-body d-flex flex-column align-items-center justify-content-center">
        <div
          class="rounded-circle bg-primary bg-opacity-10 p-3 mb-3 text-primary d-inline-flex align-items-center justify-content-center"
          style="width: 64px; height: 64px;"
        >
          <i class="bi bi-cloud-arrow-up fs-2"></i>
        </div>

        <p class="card-text fw-bold mb-1">Drag and drop your file here</p>
        <p class="card-text text-muted small mb-3">or click to browse your files</p>

        <!-- Hidden Input -->
        <input
          type="file"
          ref="fileInput"
          class="d-none"
          :disabled="isUploading"
          @change="onFileChange"
        />

        <div v-if="selectedFile" class="mt-2 w-100" @click.stop>
          <div
            class="d-inline-flex align-items-center gap-2 alert alert-primary py-2 px-3 m-0 mb-3 small text-start"
          >
            <i class="bi bi-file-earmark-check fs-5"></i>
            <span class="text-truncate" style="max-width: 250px;">{{ selectedFile.name }}</span>
            <button
              type="button"
              class="btn-close ms-auto small"
              style="font-size: 0.65rem;"
              @click="selectedFile = null"
            ></button>
          </div>

          <div>
            <button type="button" class="btn btn-primary px-4 shadow-sm" @click="startUpload">
              Upload File
            </button>
          </div>
        </div>
      </div>
    </div>
  `,
}

export { Upload }
