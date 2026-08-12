import { apiGetState, apiGetUnrImportData, apiImportMatch, apiResetState } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { generateRandomString, PAGE, theme } from "../../utils.js"
import { computed, onMounted, ref, watch } from "../../vue.js"
import { Match } from "./Match.js"
import { Upload } from "./Upload.js"

var errCnt = 0
var pollingId = null

const Import = {
  props: {},
  components: {
    Upload,
    Match,
  },
  setup: (props) => {
    const unresolved = ref([])

    const loading = ref(false)

    const refresh = async () => {
      try {
        loading.value = true
        const { data, err } = await apiGetUnrImportData()
        if (err) {
          throw err
        }
        unresolved.value = data
      } catch (error) {
        console.error("Error fetching unresolved list data:", error)
      } finally {
        loading.value = false
      }
    }

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

    const successAlert = computed(() => {
      return (
        state.value.ProcessingSrsCnt + state.value.ProcessingEpsCnt > 0 &&
        state.value.uploadError === null
      )
    })

    const pollState = async (id = generateRandomString(6)) => {
      if (pollingId === null) {
        pollingId = id
      } else if (pollingId === id) {
        console.log("Polling, id: ", id)
      } else {
        return
      }
      const { data, err } = await apiGetState()
      // if network error, then let theUI KNOW TODO
      if (err) {
        console.error(err)
        if (errCnt++ > 3) {
          state.value.uploadError = err.message
          state.value.uploadActive = false
          pollingId = null
          return
        }
      }

      let again = data.uploadActive
      // let stop = (data.uploadActive === false && state.value.uploadActive) || data.uploadError

      state.value = data

      if (again) {
        setTimeout(() => pollState(id), 2000)
      } else {
        pollingId = null
      }
    }

    const totalProgress = computed(() => {
      let num1 = state.value.UploadingSrsCnt + state.value.UploadingEpsCnt
      let deno1 = state.value.UploadingSrsCntTotal + state.value.UploadingEpsCntTotal
      let fract1 = deno1 === 0 ? 0 : (num1 / deno1) * 50
      let num2 = state.value.ProcessingSrsCnt + state.value.ProcessingEpsCnt
      let deno2 = state.value.ProcessingSrsCntTotal + state.value.ProcessingEpsCntTotal
      let fract2 = deno2 === 0 ? 0 : (num2 / deno2) * 50
      return fract1 + fract2
    })

    const entrypoint = async () => {
      pollState()
      await refresh()
    }

    onMounted(() => {
      entrypoint()
    })

    const onUploadSuccess = (data) => {
      pollState()
    }

    const onUploadStart = () => {
      state.value.uploadError = null
    }

    const handleReset = async () => {
      const { data, err } = await apiResetState()
      if (err) {
        console.error(err)
        notify(MsgType.Error, "Import", err)
        return
      }
      state.value = data
    }

    const onUploadDone = () => {}

    // const onMatchDone = async ({ TvTimeSId, MId }) => {
    //   try {
    //     const err = await apiImportMatch(TvTimeSId, MId)
    //     if (err) {
    //       throw err
    //     }
    //   } catch (error) {
    //     console.error(error)
    //     notifyError(error)
    //   } finally {
    //     refresh()
    //   }
    // }

    return {
      unresolved,
      // onMatchDone,

      loading,
      onUploadSuccess,
      onUploadStart,
      onUploadDone,
      state,
      totalProgress,
      successAlert,

      handleReset,
    }
  },
  template: /* HTML */ `
    <div class="d-flex flex-grow-1 flex-column position-relative">
      <! -- loader::: -->
      <div
        v-if="loading"
        class="position-absolute top-0 start-0 w-100 h-100 d-flex flex-column align-items-center justify-content-center bg-body bg-opacity-75 rounded"
      >
        <div
          class="spinner-border text-primary mb-3"
          role="status"
          style="width: 3rem; height: 3rem;"
        >
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div
        v-if="state.uploadActive"
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
      <! -- :::loader -->
      <h4 v-if="unresolved.Series?.length > 0 || unresolved.Episodes?.length > 0" class="mb-3">
        Unresolved - {{unresolved.Series.length}} Series, {{unresolved.Episodes.length}} Episodes
      </h4>

      <div v-if="successAlert" class="alert alert-success" role="alert">
        <i class="bi bi-check-circle-fill text-success"></i>
        Processed {{ state.ProcessingSrsCnt + state.ProcessingEpsCnt }} records
      </div>

      <div v-if="state.uploadError" class="alert alert-danger" role="alert">
        <i class="bi bi-x-circle-fill text-danger"></i>
        {{ state.uploadError }}
      </div>
      <button
        v-if="state.uploadError || successAlert"
        class="btn btn-primary align-self-end"
        role="alert"
        @click="handleReset"
      >
        <i class="bi bi-x"></i>
        Clear
      </button>
      <div v-if="state.uploadActive" class="progress" style="height: 6px;" style="z-index: 12;">
        <div
          class="progress-bar progress-bar-striped progress-bar-animated"
          :style="{ width: totalProgress + '%' }"
        ></div>
      </div>

      <Upload
        v-if="!loading"
        @success="onUploadSuccess"
        @upload-start="onUploadStart"
        @upload-done="onUploadDone"
      ></Upload>
    </div>
  `,
}
export { Import }
