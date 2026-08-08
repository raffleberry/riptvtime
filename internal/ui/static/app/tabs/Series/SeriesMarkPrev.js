import { ky, TvStatus } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, useTemplateRef, watch } from "../../vue.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { useTracked } from "../../stores/tracked.js"
import { useSeriesStore } from "./seriesStore.js"

export const SeriesMarkPrev = {
  props: {
    mid: Number,
    eps: {
      type: Array,
      default: [],
    },
    show: Boolean,
  },
  setup(props) {
    const trackedStore = useTracked()

    const { series } = storeToRefs(trackedStore)

    const { remSeries, addSeries, epMarkWatched } = useSeriesStore()

    var bSelf = null
    console.log(props.eps)

    onMounted(() => {
      const el = document.getElementById("seriesMarkPrev")
      bSelf = bootstrap.Offcanvas.getOrCreateInstance(el)
      // el.addEventListener('show.bs.offcanvas', () => {
      // })
      // el.addEventListener('hidden.bs.offcanvas', () => {
      // })
    })
    watch(
      () => props.show,
      (val) => {
        bSelf.toggle()
      },
    )

    const handleMark = async () => {
      try {
        bSelf.hide()
        await epMarkWatched(props.mid, props.eps.slice(0, 1))
      } catch (error) {
        notify(MsgType.Error, "Series", error)
      } finally {
      }
    }

    const handleMarkAll = async () => {
      try {
        bSelf.hide()
        await epMarkWatched(props.mid, props.eps)
      } catch (error) {
        notify(MsgType.Error, "Series", error)
      } finally {
      }
    }

    return {
      handleMark,
      handleMarkAll,
      ky,
    }
  },

  template: /* HTML */ `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="seriesMarkPrev">
      <div class="offcanvas-header border-bottom">
        <h5 class="offcanvas-title">Mark as watched?</h5>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="offcanvas"
          aria-label="Close"
        ></button>
      </div>
      <div class="offcanvas-body p-0">
        <div class="list-group list-group-flush">
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
            class="btn btn-primary"
            @click="handleMarkAll"
          >
            Mark all Previous too - ({{eps.length}} episodes)
          </button>
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
            class="btn btn-primary"
            @click="handleMark"
          >
            Mark as watched <span v-if="eps.length > 1"> ({{ky(eps[0].S, eps[0].E)}})</span>
          </button>
          <button
            class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-danger"
            class="btn btn-danger"
            @click="handleCancel"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  `,
}
