import { MsgType, notify } from "../../components/Notify/Notify.js"
import { ky } from "../../utils.js"
import { onMounted, ref, watch } from "../../vue.js"
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
    const { epMarkWatched } = useSeriesStore()

    const showEpsList = ref(false)

    var bSelf = null

    onMounted(() => {
      const el = document.getElementById("seriesMarkPrev")
      bSelf = bootstrap.Offcanvas.getOrCreateInstance(el)
    })
    watch(
      () => props.show,
      (val) => {
        bSelf.toggle()
        showEpsList.value = false
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

    const handleCancel = () => {
      bSelf.hide()
    }

    return {
      handleMark,
      handleMarkAll,
      handleCancel,
      ky,
      showEpsList,
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
          <!--  -->
          <div class="accordion" id="accordianEps">
            <div class="accordion-item">
              <h2 class="accordion-header">
                <button
                  class="accordion-button collapsed"
                  type="button"
                  data-bs-toggle="collapse"
                  data-bs-parent="#accordianEps"
                  data-bs-target="#accordianEpsList"
                  aria-expanded="false"
                >
                  Show {{eps.length}} episodes list
                </button>
              </h2>
              <div
                id="accordianEpsList"
                class="accordion-collapse collapse"
                :class="{ 'show': showEpsList }"
                data-bs-parent="#accordianEps"
              >
                <div class="accordion-body">
                  <ul>
                    <li v-for="ep in eps">{{ky(ep.S, ep.E)}}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
          <!--  -->

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
