import { TvStatus } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, watch } from "../../vue.js"
import { notifyError } from "../../components/Error.js"
import { useTracked } from "../../stores/tracked.js"

export const selected = ref({
  Id: null,
  Name: null,
  Year: null,
  Status: 0,
})

export const SearchTileOpts = {
  setup() {
    const loading = ref(false)

    const trackedStore = useTracked()

    const { series } = storeToRefs(trackedStore)
    const { changeStatus, remSeries, addSeries } = trackedStore

    const onDetails = () => {
      console.log("Details")
    }

    const status = computed(() => {
      let ob = series.value?.[selected.value.Id]
      if (ob) {
        return ob.TrackingStatus
      }
      return TvStatus.NotWatching
    })

    const statusActionTxt = computed(() => {
      if (status.value === TvStatus.Watching) {
        return "Stop watching"
      } else if (status.value === TvStatus.Stopped) {
        return "Resume watching"
      }
    })

    const addRemActionTxt = computed(() => {
      return status.value === TvStatus.NotWatching ? "Add Series" : "Remove Series"
    })

    var bSelf = null

    onMounted(() => {
      const el = document.getElementById("searchTileOpts")
      bSelf = new bootstrap.Offcanvas(el)

      // el.addEventListener('show.bs.offcanvas', () => {
      // })

      el.addEventListener("hidden.bs.offcanvas", () => {
        console.log("selected.value = {}", selected.value)
        selected.value = {}
      })
    })

    const handleStatusChange = async () => {
      try {
        loading.value = true

        if (status.value === TvStatus.Watching) {
          await changeStatus(selected.value.Id, TvStatus.Stopped)
        } else if (status.value === TvStatus.Stopped) {
          await changeStatus(selected.value.Id, TvStatus.Watching)
        }
        bSelf.hide()
      } catch (error) {
        notifyError(error)
      } finally {
        loading.value = false
      }
    }

    const handleAddRemSeries = async () => {
      try {
        loading.value = true
        if (status.value === TvStatus.NotWatching) {
          const err = await addSeries(selected.value.Id)
          if (err) {
            throw err
          }
        } else {
          const err = await remSeries(selected.value.Id)
          if (err) {
            throw err
          }
        }

        bSelf.hide()
      } catch (error) {
        console.error(error)
        notifyError(error)
      } finally {
        loading.value = false
      }
    }

    return {
      selected,
      loading,
      TvStatus,
      status,

      statusActionTxt,
      addRemActionTxt,

      handleStatusChange,
      handleAddRemSeries,

      onDetails,
    }
  },

  template: /* HTML */ `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="searchTileOpts">
      <div class="offcanvas-header border-bottom">
        <h5 class="offcanvas-title">
          {{selected.Name}} <span class="text-muted">({{selected.Year}})</span>
        </h5>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="offcanvas"
          aria-label="Close"
        ></button>
      </div>
      <div class="offcanvas-body p-0">
        <div v-if="loading" class="list-group list-group-flush">
          <div
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center justify-content-center border-0"
          >
            <div class="spinner-border" role="status">
              <span class="visually-hidden">Loading...</span>
            </div>
          </div>
        </div>
        <div v-else class="list-group list-group-flush">
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
          >
            Show details
          </button>
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
            v-if="[TvStatus.Watching, TvStatus.Stopped].includes(status)"
            class="btn btn-primary"
            @click="handleStatusChange"
          >
            {{ statusActionTxt }}
          </button>
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-danger"
            class="btn btn-primary"
            @click="handleAddRemSeries"
          >
            {{ addRemActionTxt }}
          </button>
        </div>
      </div>
    </div>
  `,
}
