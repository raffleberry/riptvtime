import { fmtRating, imgPosterUrl, ky } from "../../utils.js"
import { computed } from "../../vue.js"

const TvFeedTile = {
  components: {},
  props: {
    tv: Object,
  },
  emits: ["markupnext"],
  setup(props, { emit }) {
    const upNext = computed(() => ky(props.tv.UpNext.S, props.tv.UpNext.E))
    const toWatchCnt = computed(() => {
      let c = props.tv.EpisodesAired - props.tv.EpisodesWatched - 1
      return c < 0 ? 0 : c
    })
    const watchProgress = computed(() => (props.tv.EpisodesWatched / props.tv.EpisodesAired) * 100)

    const markUpNext = async () => {
      emit("markupnext", {
        MId: props.tv.MId,
        Name: props.tv.Name,
        Year: props.tv.Year,
        S: props.tv.UpNext.S,
        E: props.tv.UpNext.E,
      })
    }

    return {
      upNext,
      toWatchCnt,
      watchProgress,
      markUpNext,
      imgPosterUrl,
      fmtRating,
    }
  },

  template: /* HTML */ `
    <div class="card">
      <div class="d-flex flex-column flex-sm-row">
        <div class="tile-poster align-self-center align-self-sm-start p-2 p-sm-0">
          <img
            :src="imgPosterUrl(tv.Image)"
            class="img-fluid poster-img rounded"
            :alt="tv.Name"
            loading="lazy"
          />
        </div>

        <div class="d-flex flex-column card-body justify-content-between">
          <div>
            <h5 class="card-title d-flex flex-wrap gap-2 justify-content-between">
              <router-link :to="'/series/' + tv.MId + '#season' + tv.UpNext.S">
                {{ tv.Name }} <span class="text-muted">({{ tv.Year }})</span>
                <span class="badge bg-secondary" v-if="tv.RecentlyAired">New</span>
              </router-link>
              <span class="text-muted fs-6" v-if="tv.ImdbRating"
                >IMDb: {{ fmtRating(tv.ImdbRating) }}</span
              >
            </h5>
            <p class="card-text clamp-3">{{ tv.Overview }}</p>
          </div>
          <p class="card-text text-end mt-2">
            Up Next :
            <button
              @click="markUpNext"
              type="button"
              class="btn btn-outline-primary position-relative"
            >
              <span class="me-2"> {{ upNext }} </span>
              <span
                v-if="toWatchCnt > 0"
                class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-success"
              >
                +{{ toWatchCnt }}
                <span class="visually-hidden">unwatched episodes</span>
              </span>

              <i class="bi bi-check-circle"></i>
            </button>
          </p>
        </div>
      </div>
      <div
        class="progress rounded-top-0"
        role="progressbar"
        aria-label="Tv show progress"
        :aria-valuenow="watchProgress"
        aria-valuemin="0"
        aria-valuemax="100"
      >
        <div class="progress-bar" :style="{width: watchProgress + '%'}"></div>
      </div>
    </div>
  `,
}

export { TvFeedTile }
