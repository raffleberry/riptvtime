import { useTracked } from "../../stores/tracked.js"
import { imgPosterUrl, TvStatus } from "../../utils.js"
import { computed, ref, storeToRefs } from "../../vue.js"
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
      imgPosterUrl,
      tv,
      openSeriesOptions,
      statusCss,
    }
  },

  template: /* HTML */ `
    <div :class="statusCss.card">
      <div class="d-flex flex-column flex-sm-row">
        <div class="tile-poster align-self-center align-self-sm-start p-2 p-sm-0">
          <img
            :src="imgPosterUrl(tv.Image)"
            class="img-fluid poster-img rounded"
            :alt="tv.Name"
            loading="lazy"
          />
        </div>
        <div class="flex-grow-1" style="min-width: 0;">
          <div class="card-body">
            <div class="d-flex justify-content-between gap-2">
              <h5 class="card-title mb-0" style="min-width: 0;">
                <router-link :to="'/series/' + tv.Id">
                  {{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span>
                  <span class="badge bg-secondary"> </span>
                </router-link>
              </h5>
              <button
                @click="openSeriesOptions"
                data-bs-toggle="offcanvas"
                data-bs-target="#searchTileOpts"
                type="button"
                class="btn p-2 d-inline-flex align-items-center justify-content-center flex-shrink-0"
                aria-label="Series options"
              >
                <i :class="statusCss.icon"></i>
              </button>
            </div>
            <p class="card-text clamp-3">{{ tv.Overview }}</p>
          </div>
        </div>
      </div>
    </div>
  `,
}

export { SearchTile }
