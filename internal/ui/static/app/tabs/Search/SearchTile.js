import { useTracked } from "../../stores/tracked.js"
import { TvStatus } from "../../utils.js"
import { computed, ref, storeToRefs, watch } from "../../vue.js"
import { selected } from "./SearchTileOpts.js"

const SearchTile = {
  components: {},
  props: {
    tv: Object,
  },
  setup(props) {
    let tv = computed(() => props.tv)

    const { series } = storeToRefs(useTracked())

    const cardBorder = ref("")
    const btnIcon = ref("")

    const status = computed(() => {
      let ob = series.value?.[tv.value.Id]
      if (ob) {
        return ob.TrackingStatus
      }
      return TvStatus.NotWatching
    })

    const statusCss = computed(() => {
      let card = "card"
      let icon = "bi"
      switch (status.value) {
        case TvStatus.Watching:
          card += " border border-warning"
          icon += " bi-bookmark-check"
          break
        case TvStatus.NotWatching:
          // cb += ' border border-primary'
          icon += " bi-three-dots-vertical"
          break
        case TvStatus.Stopped:
          card += " border border-danger"
          icon += " bi-bookmark-x"
          break
        case TvStatus.Completed:
          card += " border border-success"
          icon += " bi-bookmark-check-fill"
          break
        case TvStatus.UpToDate:
          card += " border border-info"
          icon += " bi-bookmark-check-fill"
          break

        default:
          break
      }
      return {
        card,
        icon,
      }
    })

    const openSeriesOptions = () => {
      selected.value = tv.value
    }

    return {
      tv,
      openSeriesOptions,
      statusCss,
    }
  },

  template: /* HTML */ `
    <div :class="statusCss.card">
        <div class="row">
            <!--
            <div class="col-md-4">
            <img src="..." class="img-fluid rounded-start" alt="..."> 
            </div>
            <div class="col-md-8">
            </div>
            -->
            <div class="col">
                <div class="card-body">
                    <div class="d-flex justify-content-between">
                        <h5 class="card-title">
                            <router-link :to="'/series/' + tv.Id">
                            {{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span> <span class="badge bg-secondary" </span>
                            </router-link>
                        </h5>
                        <button @click="openSeriesOptions" 
                            data-bs-toggle="offcanvas" data-bs-target="#searchTileOpts"
                            type="button" class="btn p-2 d-inline-flex align-items-center justify-content-center">
                                <i :class="statusCss.icon"></i>
                        </button>
                    </div>
                    <p class="card-text">{{ tv.Overview }}</p>
                </div>
            </div>
        </div>
    </div>
`,
}

export { SearchTile }
