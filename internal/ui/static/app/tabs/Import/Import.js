import { apiGetState, apiGetUnMatchedImportList, apiImportMatch } from "../../api.js"
import { notifyError } from "../../components/Error.js"
import { PAGE, theme } from "../../utils.js"
import { computed, onMounted, ref, watch } from "../../vue.js"
import { Match } from "./Match.js"
import { Upload } from "./Upload.js"

const Import = {
  props: {},
  components: {
    Upload,
    Match,
  },
  setup: (props) => {
    const Mloading = ref(false)

    const unMatched = ref([])

    const refresh = async () => {
      try {
        Mloading.value = true
        const { data, err } = await apiGetUnMatchedImportList()
        if (err) {
          throw err
        }
        unMatched.value = data
      } catch (error) {
        console.error("Error fetching unmatched list data:", error)
      } finally {
        Mloading.value = false
      }
    }

    const Uloading = ref(false)

    const pollStateId = ref(null)

    let errCnt = 0

    const successAlert = ref(false)

    const state = ref({
      IUploadingSrsCnt: 0,
      IUploadingSrsCntTotal: 0,
      IUploadingEpsCnt: 0,
      IUploadingEpsCntTotal: 0,
      IUploadActive: false,
      IUploadError: null,
    })

    watch(state, (nv, ov) => {
      if (!nv.IUploadActive && ov.IUploadActive) {
        if (!nv.IUploadError) {
          successAlert.value = true
        }
        clearInterval(pollStateId.value)
        pollStateId.value = null
      }
    })

    const IUploadProgress = computed(() => {
      return (
        ((state.value.IUploadingSrsCnt + state.value.IUploadingEpsCnt) /
          (state.value.IUploadingSrsCntTotal + state.value.IUploadingEpsCntTotal)) *
        100
      )
    })

    onMounted(() => {
      refresh()
    })

    const onUploadSuccess = (data) => {
      // update upload count
      refresh()
    }

    const getState = async () => {
      const { data, err } = await apiGetState()
      if (err) {
        if (errCnt++ > 3) {
          clearInterval(pollStateId.value)
          pollStateId.value = null
          state.value.IUploadError = err.message
          state.value.IUploadActive = false
          Uloading.value = false
        }
        console.log(err)
        return
      }
      state.value = data
    }

    const onUploadStart = () => {
      successAlert.value = false
      state.value.IUploadError = null
      Uloading.value = true
    }

    const onUploadDone = () => {
      Uloading.value = false
      if (pollStateId.value) {
        console.warn("POLLING ALREADY ACTIVE!!")
        clearInterval(pollStateId.value)
      }
      pollStateId.value = setInterval(getState, 2000)
    }

    const onMatchDone = async ({ TvTimeSId, MId }) => {
      try {
        Mloading.value = true
        const err = await apiImportMatch(TvTimeSId, MId)
        if (err) {
          throw err
        }
      } catch (error) {
        console.error(error)
        notifyError(error)
      } finally {
        Mloading.value = false
        refresh()
      }
    }

    return {
      unMatched,
      onMatchDone,
      Mloading,

      loading: Uloading,
      onUploadSuccess,
      onUploadStart,
      onUploadDone,
      state,
      pollStateId,
      IUploadProgress,
      successAlert,
    }
  },
  template: /* HTML */ `
    <div class="container mt-4">
      <div class="position-relative">
        <div
          v-if="Mloading"
          class="mx-3 position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-25 rounded"
          style="z-index: 11; backdrop-filter: blur(3px); -webkit-backdrop-filter: blur(3px);"
        >
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <Match v-if="unMatched?.length > 0" :series="unMatched" @match-done="onMatchDone"></Match>
      </div>

      <div v-if="successAlert" class="alert alert-success" role="alert">
        <i class="bi bi-check-circle-fill text-success"></i>
        Processed {{ state.IUploadingSrsCnt + state.IUploadingEpsCnt }} records
      </div>

      <div v-if="state.IUploadError" class="alert alert-danger" role="alert">
        <i class="bi bi-x-circle-fill text-danger"></i>
        {{ state.IUploadError }}
      </div>
      <div v-if="pollStateId" class="progress" style="height: 6px;" style="z-index: 12;">
        <div
          class="progress-bar progress-bar-striped progress-bar-animated"
          :style="{ width: IUploadProgress + '%' }"
        ></div>
      </div>

      <div class="position-relative">
        <div
          v-if="pollStateId"
          class="position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-75 rounded"
          style="z-index: 11; backdrop-filter: blur(3px); -webkit-backdrop-filter: blur(3px);"
        >
          <p class="text-center">
            Series: {{ state.IUploadingSrsCnt }} / {{ state.IUploadingSrsCntTotal }}
          </p>
          <p class="text-center">
            Episodes: {{ state.IUploadingEpsCnt }} / {{ state.IUploadingEpsCntTotal }}
          </p>
        </div>

        <Upload
          @success="onUploadSuccess"
          @upload-start="onUploadStart"
          @upload-done="onUploadDone"
        ></Upload>
      </div>
    </div>
  `,
}
export { Import }
