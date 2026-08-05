import { apiGetState, apiGetUnMatchedImportData, apiImportMatch } from "../../api.js"
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
      return // TODO
      try {
        Mloading.value = true
        const { data, err } = await apiGetUnMatchedImportData()
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
      UploadingSrsCnt: 0,
      UploadingSrsCntTotal: 0,

      UploadingEpsCnt: 0,
      UploadingEpsCntTotal: 0,

      ProcessingSrsCnt: 0,
      ProcessingSrsCntTotal: 0,

      ProcessingEpsCnt: 0,
      ProcessingEpsCntTotal: 0,

      Stage: 0,
      StageCnt: 2,

      uploadActive: false,
      uploadError: null,
    })

    watch(state, (nv, ov) => {
      if (!nv.uploadActive && ov.uploadActive) {
        if (!nv.uploadError) {
          successAlert.value = true
        }
        clearInterval(pollStateId.value)
        pollStateId.value = null
      }
    })

    const totalProgress = computed(() => {
      const num =
        state.value.UploadingSrsCnt +
        state.value.UploadingEpsCnt +
        state.value.ProcessingSrsCnt +
        state.value.ProcessingEpsCnt
      const deno =
        state.value.UploadingSrsCntTotal +
        state.value.UploadingEpsCntTotal +
        state.value.ProcessingSrsCntTotal +
        state.value.ProcessingEpsCntTotal
      if (deno == 0) return 0
      return (num / deno) * 100
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
      // if network error, then let theUI KNOW
      if (err) {
        if (errCnt++ > 3) {
          clearInterval(pollStateId.value)
          pollStateId.value = null
          state.value.uploadError = err.message
          state.value.uploadActive = false
          Uloading.value = false
        }
        console.log(err)
        return
      }
      state.value = data
    }

    const onUploadStart = () => {
      successAlert.value = false
      state.value.uploadError = null
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
    window.poll = onUploadDone

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
      totalProgress,
      successAlert,
    }
  },
  template: /* HTML */ `
    <div class="mt-4">
      <div class="position-relative">
        <div
          v-if="Mloading"
          class="position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-25 rounded"
          style="z-index: 11; backdrop-filter: blur(3px); -webkit-backdrop-filter: blur(3px);"
        >
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <!-- TODO
        <Match
          v-if="unMatched.Series?.length > 0 || unMatched.Episodes?.length > 0"
          :data="unMatched"
          @match-done="onMatchDone"
        ></Match>
        -->
      </div>

      <div v-if="successAlert" class="alert alert-success" role="alert">
        <i class="bi bi-check-circle-fill text-success"></i>
        Processed {{ state.ProcessingSrsCnt + state.ProcessingEpsCnt }} records
      </div>

      <div v-if="state.uploadError" class="alert alert-danger" role="alert">
        <i class="bi bi-x-circle-fill text-danger"></i>
        {{ state.uploadError }}
      </div>
      <div v-if="pollStateId" class="progress" style="height: 6px;" style="z-index: 12;">
        <div
          class="progress-bar progress-bar-striped progress-bar-animated"
          :style="{ width: totalProgress + '%' }"
        ></div>
      </div>

      <div class="position-relative">
        <div
          v-if="pollStateId"
          class="position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-75 rounded"
          style="z-index: 11; backdrop-filter: blur(3px); -webkit-backdrop-filter: blur(3px);"
        >
          <div>
            <p class="text-center">Stage: {{ state.Stage }} / {{ state.StageCnt }}</p>
            <p v-if="state.Stage == 1" class="text-center">
              Series: {{ state.UploadingSrsCnt }} / {{ state.UploadingSrsCntTotal }}
            </p>
            <p v-if="state.Stage == 1" class="text-center">
              Episodes: {{ state.UploadingEpsCnt }} / {{ state.UploadingEpsCntTotal }}
            </p>
            <p v-if="state.Stage == 2" class="text-center">
              Series: {{ state.ProcessingSrsCnt }} / {{ state.ProcessingSrsCntTotal }}
            </p>
            <p v-if="state.Stage == 2" class="text-center">
              Episodes: {{ state.ProcessingEpsCnt }} / {{ state.ProcessingEpsCntTotal }}
            </p>
          </div>
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
